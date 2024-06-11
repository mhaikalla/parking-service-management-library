package jwt

import (
	"encoding/json"
	"errors"

	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jws"
)

// JWT jwt context consisting config
type JWT struct {
	isJWE               bool
	key                 interface{}
	signingMethod       jwa.SignatureAlgorithm
	encryptionMethod    jwa.ContentEncryptionAlgorithm
	compressionMethod   jwa.CompressionAlgorithm
	keyEncryptionMethod jwa.KeyEncryptionAlgorithm
	// serializer          func(string) (JWTClaims, error)
	assertFn func(JWTClaims) error
	enable   bool
	duration int
}

func (j *JWT) serializeClaims(claims interface{}) (string, error) {
	buff, err := json.Marshal(claims)

	if err != nil {
		return "", err
	}

	sig, err := jws.Sign(buff, j.signingMethod, j.key)

	if err != nil {
		return "", err
	}
	return string(sig), nil
}

// WithMapClaims giving map[string]interface{} as map claims return string token or an error
func (j *JWT) WithMapClaims(claims map[string]interface{}) (string, error) {
	claims["exp"] = time.Now().Add(time.Second * time.Duration(j.duration)).Unix()
	return j.serializeClaims(claims)
}

// WithStructClaims giving JWTClaims as claims return string token or an error
func (j *JWT) WithStructClaims(claims JWTClaims) (string, error) {
	claims.ExpiredAt = time.Now().Add(time.Second * time.Duration(j.duration)).Unix()
	return j.serializeClaims(claims)
}

// Duration return duration in second this jwt settled
func (j *JWT) Duration() int {
	return j.duration
}

// NewClaims ...
func (j *JWT) NewClaims() JWTClaims {
	return JWTClaims{
		JWTID:    uuid.New().String(),
		IssuedAt: time.Now().Unix(),
	}
}

// Serialize attempt to serialize compact jwt string to JWTClaims
func (j *JWT) Serialize(compactedJWT string) (JWTClaims, error) {
	res := JWTClaims{}
	var buff []byte
	var err error

	buff, err = verifyWithKey([]byte(compactedJWT), j.signingMethod, j.key)

	if err := json.Unmarshal(buff, &res); err != nil {
		return res, err
	}

	return res, err
}

// Bind binding jwt to custom claims
func (j *JWT) Bind(compactedJWT string, customClaims interface{}) error {
	if j.isJWE {
		buff, err := jwe.Decrypt([]byte(compactedJWT), j.keyEncryptionMethod, j.key)
		if err != nil {
			return err
		}
		return json.Unmarshal(buff, &customClaims)
	}
	buff, err := verifyWithKey([]byte(compactedJWT), j.signingMethod, j.key)
	if err != nil {
		return err
	}
	return json.Unmarshal(buff, &customClaims)
}

// GetRaw get raw JWT in bytes
func (j *JWT) GetRaw(compactedJWT string) ([]byte, error) {
	return verifyWithKey([]byte(compactedJWT), j.signingMethod, j.key)
}

// Assert take JWTClaims as claims and error if fail
func (j *JWT) Assert(claims JWTClaims) error {
	return j.assertFn(claims)
}

// NewJWT create a new jwt context
func NewJWT(config map[string]map[string]interface{}) IJWT {

	var ctx JWT

	if conf, found := config["jwt"]; found {

		if secConf, found := config["secrets"]; found {
			if secJWTConf, found := secConf["jwt"].(map[string]interface{}); found {

				if _, found := conf["encryption_method"]; found {
					ctx = parseJWEConfig(conf, secJWTConf)
				} else {
					ctx = parseJWSConfig(conf, secJWTConf)
				}

			} else {
				panic(errors.New("No configuration key jwt found in secrets entry"))
			}
		} else {
			panic(errors.New("No configuration key secrets found"))
		}
	} else {
		panic(errors.New("No configuration key jwt found"))
	}

	return &ctx
}
