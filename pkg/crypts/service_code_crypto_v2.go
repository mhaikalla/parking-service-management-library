package crypts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"
)

// Contains env var to enable crypto.
var (
	ServiceCodeCryptoEnabled = os.Getenv("SERVICE_CODE_CRYPTO")
)

func unpackNonEncrypted(compactCode string) (sid, ofid, channel, location string) {
	if !strings.Contains(compactCode, "|") {
		return compactCode, "", "", ""
	}
	splited := strings.Split(compactCode, "|")
	if len(splited) == 3 {
		return splited[0], splited[1], splited[2], ""
	}
	return splited[0], splited[1], splited[2], splited[3]
}

func unpackNonEncryptedWithTime(compactCode string) (sid, ofid, channel, location string, timestamp time.Time) {
	if !strings.Contains(compactCode, "|") {
		return compactCode, "", "", "", time.Now().Local()
	}
	splited := strings.Split(compactCode, "|")
	lenSplitted := len(splited)
	switch lenSplitted {
	case 2:
		timeStamp, parseErr := time.Parse(time.RFC3339, splited[1])
		panicWhenErr(parseErr)
		return splited[0], "", "", "", timeStamp
	case 3:
		timeStamp, parseErr := time.Parse(time.RFC3339, splited[2])
		panicWhenErr(parseErr)
		return splited[0], splited[1], "", "", timeStamp
	case 4:
		timeStamp, parseErr := time.Parse(time.RFC3339, splited[3])
		panicWhenErr(parseErr)
		return splited[0], splited[1], splited[2], "", timeStamp
	case 5:
		timeStamp, parseErr := time.Parse(time.RFC3339, splited[4])
		panicWhenErr(parseErr)
		return splited[0], splited[1], splited[2], splited[3], timeStamp
	default:
		return compactCode, "", "", "", time.Now().Local()
	}

}

func packNonEncrypted(code ...string) string {
	switch {
	case len(code) < 1:
		return ""
	case len(code) == 1:
		return code[0]
	default:
		return strings.Join(code, "|")
	}
}

func panicWhenErr(err error) {
	if err != nil {
		panic(err)
	}
}

// =====  Service Code Object =====

// ServiceCodeV2 like `ServiceCode`.
type ServiceCodeV2 struct {
	serviceID string
	offerID   string
	channel   string
	location  string
	timestamp time.Time
}

// GetServiceCode get Service ID from Unpacked version.
func (sc ServiceCodeV2) GetServiceCode() string {
	return sc.serviceID
}

// GetOfferCode get Offer ID from Unpacked version.
func (sc ServiceCodeV2) GetOfferCode() string {
	return sc.offerID
}

// GetChannel get Offer Channel from Unpacked version.
func (sc ServiceCodeV2) GetChannel() string {
	return sc.channel
}

// GetTimestamp get timestamp generated when unpacking.
func (sc ServiceCodeV2) GetTimestamp() time.Time {
	return sc.timestamp
}

// IsExpired check if timestamp is expired, see `crypto#ServiceCodeCryptoV2.GetTimestamp`.
func (sc ServiceCodeV2) IsExpired() bool {
	currentTime := time.Now().Local()
	return currentTime.After(sc.timestamp.Add(time.Duration(expirySess) * time.Second))
}

// IsOffer check if this service code contains identifier for MCCM offer,
// see also `crypto#ServiceCodeCryptoV2.GetOfferCode` and `crypto#ServiceCodeCryptoV2.GetChannel`
func (sc ServiceCodeV2) IsOffer() bool {
	return sc.offerID != "" && sc.channel != ""
}

// IsLocationExist check if this service code contains page location of MCCM offer.
func (sc ServiceCodeV2) IsLocationExist() bool {
	return sc.location != ""
}

// GetLocation get page location of this service code if exists.
func (sc ServiceCodeV2) GetLocation() string {
	return sc.location
}

// =====  Service Code Crypto Object =====

// ServiceCodeCryptoV2 contains data needed to packing or unpackage service code.
type ServiceCodeCryptoV2 struct {
	msisdn    string
	subsid    string
	substype  string
	isEnabled bool
	validator func(time.Time) bool
}

// UnPack unpack service code, explode all element it contains to use later,
// panic when error.
func (scc ServiceCodeCryptoV2) UnPack(compactCode string) ServiceCodeV2 {
	res := ServiceCodeV2{}

	if !scc.isEnabled {
		res.serviceID, res.offerID, res.channel, res.location = unpackNonEncrypted(compactCode)
		res.timestamp = time.Now().Local()
		return res
	}

	buff, base64Err := base64.RawURLEncoding.DecodeString(compactCode)
	panicWhenErr(base64Err)

	buff = buff[len(serviceCodePrefix):]
	derivedKey := saltKey(nounceKey, scc.msisdn, scc.subsid)
	block, cipherErr := aes.NewCipher(derivedKey)
	panicWhenErr(cipherErr)

	gcm, gcmErr := cipher.NewGCM(block)
	panicWhenErr(gcmErr)

	nonceSize := gcm.NonceSize()
	nonce, data := buff[:nonceSize], buff[nonceSize:]
	dst := make([]byte, 0)

	dec, decErr := gcm.Open(dst, nonce, data, []byte(scc.substype))
	panicWhenErr(decErr)

	// splited := strings.Split(string(dec), "|")
	a, b, c, d, e := unpackNonEncryptedWithTime(string(dec))
	res.serviceID, res.offerID, res.channel, res.location, res.timestamp = a, b, c, d, e

	if !scc.validator(e) {
		panicWhenErr(errors.New(sessionExpiredError))
	}

	return res
}

// Pack packing all args in to service code,
// could be one string with layout `serviceID|offerID|Channel|pageLocations`,
// or one element per args with the same sequence.
func (scc ServiceCodeCryptoV2) Pack(code ...string) string {
	plainText := packNonEncrypted(code...)
	timestampStr := time.Now().Local().Format(time.RFC3339)

	if !scc.isEnabled {
		return plainText
	}

	derivedKey := saltKey(nounceKey, scc.msisdn, scc.subsid)
	block, cipherErr := aes.NewCipher(derivedKey)
	panicWhenErr(cipherErr)

	gcm, gcmErr := cipher.NewGCM(block)
	panicWhenErr(gcmErr)

	nonce := make([]byte, gcm.NonceSize())
	dst := []byte(serviceCodePrefix)

	if _, randErr := rand.Read(nonce); randErr != nil {
		panicWhenErr(randErr)
	}

	plainText += "|" + timestampStr
	res := gcm.Seal(nonce, nonce, []byte(plainText), []byte(scc.substype))
	dst = append(dst, res...)
	return base64.RawURLEncoding.EncodeToString(dst)
}

// NewServiceCodeCryptoV2 create a new `crypto#ServiceCodeCryptoV2` using `ctx`,
// `ctx` is equivalent with `contexts.BearerContext`.
func NewServiceCodeCryptoV2(ctx IInterceptorContext) ServiceCodeCryptoV2 {
	validator := func(time.Time) bool { return true }

	if len(strictLocations) > 0 {
		callerFunc := getCallerFunc()
		if found := strictLocations[callerFunc]; found {
			validator = createExpiryValidator(expirySess)
		}
	}
	return ServiceCodeCryptoV2{
		isEnabled: ServiceCodeCryptoEnabled == "1",
		validator: validator,
		subsid:    ctx.GetSubsID(),
		msisdn:    ctx.GetMSISDN(),
		substype:  ctx.GetSubsType(),
	}
}
