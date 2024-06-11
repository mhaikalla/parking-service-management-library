package crypts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"parking-service/pkg/condutils"

	"golang.org/x/crypto/pbkdf2"
)

const (
	serviceCodePrefix   = `SC__`
	sessionExpiredError = `your service code session expired, please refresh it`
	// initDebugMessage    = `running service code crypto module with cipher length: %v, strict locations: %v, additional salt: %v, expiry session: %v sec`
)

var (
	nounceKey       = condutils.Or(os.Getenv("ITEM_SALT_KEY"), "7Bu3hd69eZx5X0jJEDzuNT54uK46md").(string)
	strictLocations = map[string]bool{}
	expirySess      = 60
)

// ServiceCode contains all data related to service code.
type ServiceCode struct {
	serviceID string
	offerID   string
	channel   string
	timestamp time.Time
}

// GetServiceCode get Service ID from Unpacked version.
func (sc ServiceCode) GetServiceCode() string {
	return sc.serviceID
}

// GetOfferCode get Offer ID from Unpacked version.
func (sc ServiceCode) GetOfferCode() string {
	return sc.offerID
}

// GetChannel get Offer Channel from Unpacked version.
func (sc ServiceCode) GetChannel() string {
	return sc.channel
}

// GetTimestamp get timestamp generated when unpacking.
func (sc ServiceCode) GetTimestamp() time.Time {
	return sc.timestamp
}

// IsExpired check if timestamp is expired, see `crypto#ServiceCodeCrypto.GetTimestamp`.
func (sc ServiceCode) IsExpired() bool {
	currentTime := time.Now().Local()
	return currentTime.After(sc.timestamp.Add(time.Duration(expirySess) * time.Second))
}

// IsOffer check if this service code contains identifier for MCCM offer,
// see also `crypto#ServiceCodeCrypto.GetOfferCode` and `crypto#ServiceCodeCrypto.GetChannel`
func (sc ServiceCode) IsOffer() bool {
	return sc.offerID != "" && sc.channel != ""
}

func init() {
	sccStrictLocs := os.Getenv("SCC_STRICT_LOCS")
	sccExp := os.Getenv("SCC_EXP_SEC")

	if sccStrictLocs != "" {
		locs := strings.Split(sccStrictLocs, ";")
		for _, loc := range locs {
			strictLocations[loc] = true
		}
	}

	if sccExp != "" && regexp.MustCompile(`\d+`).MatchString(sccExp) {
		parseInt, parseErr := strconv.Atoi(sccExp)
		condutils.When(parseErr != nil, log.Println, parseErr)
		expirySess = condutils.Or(parseInt, expirySess).(int)
	}
}

// IInterceptorContext interface to get all data needed to create `ServiceCodeCrypto`.
type IInterceptorContext interface {
	GetMSISDN() string
	GetSubsID() string
	GetSubsType() string
}

// ServiceCodeCrypto contains data needed to packing or unpackage service code.
type ServiceCodeCrypto struct {
	context   IInterceptorContext
	validator func(time.Time) bool
}

func saltKey(keys ...string) []byte {
	if len(keys) < 2 {
		panic(errors.New("at least 2 arg must be supplied"))
	}

	var key, salt []byte

	if len(keys) == 2 {
		key = []byte(keys[0])
		salt = []byte(keys[1])
	} else {
		keysLen := len(keys) - 1
		key = []byte(strings.Join(keys[:keysLen], ""))
		salt = []byte(keys[keysLen])
	}

	lenChar := len(strings.Join(keys, ""))
	return pbkdf2.Key(key, salt, lenChar, aes.BlockSize, sha256.New)
}

func getCallerFunc() string {

	pc, file, _, ok := runtime.Caller(2)
	if !ok {
		panic(file)
	}

	return runtime.FuncForPC(pc).Name()
}

// UnPack unpack service code, explode all element it contains to use later.
func (scc ServiceCodeCrypto) UnPack(compactCode string) (ServiceCode, error) {
	res := ServiceCode{}

	buff, base64Err := base64.RawURLEncoding.DecodeString(compactCode)
	if base64Err != nil {
		return res, base64Err
	}

	buff = buff[len(serviceCodePrefix):]

	derivedKey := saltKey(nounceKey, scc.context.GetMSISDN(), scc.context.GetSubsID())

	block, cipherErr := aes.NewCipher(derivedKey)

	if cipherErr != nil {
		return res, cipherErr
	}

	gcm, gcmErr := cipher.NewGCM(block)

	if gcmErr != nil {
		return res, gcmErr
	}
	nonceSize := gcm.NonceSize()
	nonce, data := buff[:nonceSize], buff[nonceSize:]
	dst := make([]byte, 0)

	dec, decErr := gcm.Open(dst, nonce, data, []byte(scc.context.GetSubsType()))
	if decErr != nil {
		return res, decErr
	}

	splited := strings.Split(string(dec), "|")
	res.serviceID = splited[0]
	res.offerID = splited[1]
	res.channel = splited[2]

	parsed, parseErr := time.Parse(time.RFC3339, splited[3])
	if parseErr != nil {
		return res, parseErr
	}

	if !scc.validator(parsed) {
		return res, errors.New(sessionExpiredError)
	}

	res.timestamp = parsed
	return res, nil
}

// Pack packing all args in to service code,
// one element per args with the layout`serviceID|offerID|Channel`.
func (scc ServiceCodeCrypto) Pack(serviceID, offerID, channel string) (string, error) {
	timestampStr := time.Now().Local().Format(time.RFC3339)

	derivedKey := saltKey(nounceKey, scc.context.GetMSISDN(), scc.context.GetSubsID())

	block, cipherErr := aes.NewCipher(derivedKey)

	if cipherErr != nil {
		return "", cipherErr
	}

	gcm, gcmErr := cipher.NewGCM(block)

	if gcmErr != nil {
		return "", gcmErr
	}

	nonce := make([]byte, gcm.NonceSize())
	dst := []byte(serviceCodePrefix)

	if _, randErr := rand.Read(nonce); randErr != nil {
		return "", randErr
	}

	plainText := serviceID + "|" + offerID + "|" + channel + "|" + timestampStr

	res := gcm.Seal(nonce, nonce, []byte(plainText), []byte(scc.context.GetSubsType()))

	dst = append(dst, res...)

	return base64.RawURLEncoding.EncodeToString(dst), nil
}

func createExpiryValidator(expiry int) func(time.Time) bool {
	return func(timestamp time.Time) bool {
		currentTime := time.Now().Local()
		return !currentTime.After(timestamp.Add(time.Duration(expiry) * time.Second))
	}
}

// NewServiceCodeCrypto create a new `crypto#ServiceCodeCrypto` using `ctx`,
// `ctx` is equivalent with `contexts.BearerContext`.
func NewServiceCodeCrypto(ctx IInterceptorContext) ServiceCodeCrypto {
	validator := func(time.Time) bool { return true }

	if len(strictLocations) > 0 {
		callerFunc := getCallerFunc()
		if found := strictLocations[callerFunc]; found {
			validator = createExpiryValidator(expirySess)
		}
	}
	return ServiceCodeCrypto{context: ctx, validator: validator}
}
