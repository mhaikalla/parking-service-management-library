package jwt

import (
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/mhaikalla/parking-service-management-library/pkg/condutils"

	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"

	"github.com/lestrrat-go/jwx/jws"
	"golang.org/x/crypto/pbkdf2"
)

var CiamEnabledStatus = os.Getenv("CIAM_ENABLED_STATUS")

// deriveKey ...
func deriveKey(key string, salt interface{}, length int) []byte {
	if condutils.IsEmpty(salt) {
		return pbkdf2.Key([]byte(key), []byte(""), 100, length, sha512.New)
	}
	return pbkdf2.Key([]byte(key), []byte(salt.(string)), 100, length, sha512.New)
}

// isPrivateKey ...
func isPrivateKey(privKey string) bool {
	return strings.Contains(privKey, "PRIVATE KEY")
}

// createPrivateKeyParser ...
func createPrivateKeyParser(privateKey string) func(buff []byte) (privKey interface{}, err error) {
	if strings.Contains(privateKey, "BEGIN RSA PRIVATE KEY") {
		return func(b []byte) (interface{}, error) { return x509.ParsePKCS1PrivateKey(b) }
	}
	if strings.Contains(privateKey, "BEGIN PRIVATE KEY") {
		return x509.ParsePKCS8PrivateKey
	}
	return nil
}

// parseRSAPrivKeyString ...
func parseRSAPrivKeyString(privateKey string) interface{} {
	block, _ := pem.Decode([]byte(privateKey))
	parser := createPrivateKeyParser(privateKey)
	if parser == nil {
		panic("we not supported private key you provided, we support unencrypted RSA #PCKS1 and #PKCS8 private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return priv
}

// verifyWithKey verify with key, the key maybe private key
func verifyWithKey(buff []byte, signMethod jwa.SignatureAlgorithm, key interface{}) ([]byte, error) {
	if privKey, ok := key.(*rsa.PrivateKey); ok {
		return jws.Verify(buff, signMethod, privKey.PublicKey)
	}
	return jws.Verify(buff, signMethod, key)
}

// defaultAssertFn ...
func defaultAssertFn(claims JWTClaims) error {

	expiry := time.Unix(claims.ExpiredAt, 0)
	isCiamEnabled, _ := strconv.ParseBool(CiamEnabledStatus)

	if time.Now().Add(time.Second).After(expiry) && isCiamEnabled {
		return errors.New("token is expired")
	}

	if claims.ClientID == "" {
		return errors.New("missing Client ID")
	}

	if claims.Subject == "" {
		return errors.New("missing Subject")
	}

	if claims.Audience == "" {
		return errors.New("missing Audience")
	}

	if claims.JWTID == "" {
		return errors.New("missing JWT ID")
	}

	if claims.SessionID == "" {
		return errors.New("missing Session ID")
	}

	if claims.IssuedAt == 0 {
		return errors.New("missing Issued At")
	}

	return nil
}

// parseJWEConfig ...
func parseJWEConfig(config, secret map[string]interface{}) JWT {
	res := JWT{}

	encryptionMethod := config["encryption_method"].(string)
	compressionMethod := config["compression_method"].(string)
	key := secret["key"].(string)
	salt := secret["salt"]

	var encMethod jwa.ContentEncryptionAlgorithm
	var keyEncMethod jwa.KeyEncryptionAlgorithm
	var derivedKey interface{}

	if !condutils.IsEmpty(config["key_algo"]) {
		encMethod = jwa.ContentEncryptionAlgorithm(encryptionMethod)
		keyEncMethod = jwa.KeyEncryptionAlgorithm(config["key_algo"].(string))
		derivedKey = parseRSAPrivKeyString(key)
	} else {
		encMethod, keyEncMethod, derivedKey = EncryptionMethodPairs(encryptionMethod, key, salt)
	}

	res.isJWE = true
	res.encryptionMethod = encMethod
	res.keyEncryptionMethod = keyEncMethod
	res.key = derivedKey
	res.compressionMethod = CompressionMethodSelection(compressionMethod)
	res.enable = config["enable"].(bool)
	res.duration = config["duration"].(int)
	res.assertFn = defaultAssertFn

	return res
}

// parseJWSConfig ...
func parseJWSConfig(config, secret map[string]interface{}) JWT {
	res := JWT{}

	signingMethod := config["signing_method"].(string)
	key := secret["key"].(string)

	if isPrivateKey(key) {
		res.key = parseRSAPrivKeyString(key)
	} else {
		res.key = []byte(key)
	}

	res.isJWE = false
	res.signingMethod = SigningMethodSelection(signingMethod)
	res.enable = config["enable"].(bool)
	res.duration = config["duration"].(int)
	res.assertFn = defaultAssertFn

	return res
}

// SigningMethodSelection ...
func SigningMethodSelection(method string) jwa.SignatureAlgorithm {
	switch method {
	case "HS384":
		return jwa.HS384
	case "HS512":
		return jwa.HS512
	case "ES256":
		return jwa.ES256
	case "ES384":
		return jwa.ES384
	case "ES512":
		return jwa.ES512
	case "RS256":
		return jwa.RS256
	case "RS384":
		return jwa.RS384
	case "RS512":
		return jwa.RS512
	default:
		return jwa.HS256
	}
}

// EncryptionMethodPairs ...
func EncryptionMethodPairs(method, key string, salt interface{}) (jwa.ContentEncryptionAlgorithm, jwa.KeyEncryptionAlgorithm, interface{}) {
	switch method {
	case "A192CBC-HS384":
		return jwa.A192CBC_HS384, jwa.A192KW, deriveKey(key, salt, 48)
	case "A256CBC-HS512":
		return jwa.A256CBC_HS512, jwa.A256KW, deriveKey(key, salt, 64)
	case "A128GCM":
		return jwa.A128GCM, jwa.A128KW, deriveKey(key, salt, 32)
	case "A192GCM":
		return jwa.A192GCM, jwa.A192KW, deriveKey(key, salt, 48)
	case "A256GCM":
		return jwa.A256GCM, jwa.A256KW, deriveKey(key, salt, 64)
	default:
		return jwa.A128CBC_HS256, jwa.A128KW, deriveKey(key, salt, 32)
	}
}

// CompressionMethodSelection ...
func CompressionMethodSelection(method string) jwa.CompressionAlgorithm {
	switch method {
	case "deflate":
		return jwa.Deflate
	default:
		return jwa.NoCompress
	}
}

func errorWhenFieldEmpty(fieldName, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf(FieldStringCheckErrorMsg, fieldName, value)
	}
	return nil
}

func validateCIAMClaims(claims CIAMClaims) error {
	subsType := strings.TrimSpace(condutils.Or(claims.SubsType, claims.SubscriberType).(string))
	msisdn := claims.MSISDN
	nonFiberSubsID := claims.SubscriberID
	deviceID := claims.DeviceID
	accountID := claims.AccountID
	custID := claims.CustomerID
	email := claims.Email

	var theErrors interface{}

	switch strings.ToUpper(subsType) {
	case "":
		return fmt.Errorf(UnknownSubsType, subsType, "subscription_type", "subscriber-type")

	case "HOMEFIBER", "HOMESATU":
		theErrors = condutils.Ors(
			errorWhenFieldEmpty("accountID", accountID),
			errorWhenFieldEmpty("customerID", custID),
			errorWhenFieldEmpty("email", email),
		)

	default:
		theErrors = condutils.Ors(
			errorWhenFieldEmpty("msisdn", msisdn),
			errorWhenFieldEmpty("subscriberID", nonFiberSubsID),
			errorWhenFieldEmpty("deviceID", deviceID),
		)
	}

	if !condutils.IsEmpty(theErrors) {
		return theErrors.(error)
	}
	return nil
}

// TransCIAMClaims transform `model.CIAMClaims` to `model.JWTClaims`.
func TransCIAMClaims(claims CIAMClaims) (JWTClaims, error) {
	cClaims := JWTClaims{}

	if validateErr := validateCIAMClaims(claims); validateErr != nil {
		return cClaims, validateErr
	}

	subsType := strings.TrimSpace(condutils.Or(claims.SubsType, claims.SubscriberType).(string))

	cClaims.JWTID = claims.Subject
	cClaims.Issuer = claims.Issuer
	cClaims.ClientID = claims.DeviceID + ";" + claims.SubscriberID + ";" + claims.MSISDN
	cClaims.Subject = claims.Audience
	cClaims.Audience = claims.SubsType
	cClaims.SessionID = claims.AuditTrackingID
	cClaims.IssuedAt = claims.IssuedAt
	cClaims.ExpiredAt = claims.ExpiredAt
	cClaims.MSISDN = claims.MSISDN
	cClaims.SubsID = claims.SubscriberID
	cClaims.DeviceID = claims.DeviceID
	cClaims.IsFirstLogin = claims.IsFirstLogin

	if subsType == HomeFiberSubsType || subsType == HomeSatuSubsType {
		subId := condutils.Or(claims.SubscriberID, claims.AccountID).(string)
		subType := condutils.Or(claims.SubscriberType, claims.SubsType).(string)
		cClaims.ClientID = ";" + subId + ";" + claims.AccountID
		cClaims.Audience = subType
		cClaims.SubsID = subId
		cClaims.MSISDN = claims.AccountID
	}

	return cClaims, nil
}
