package payloadhooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"parking-service/pkg/condutils"
	"parking-service/pkg/crypts"
)

const (
	CipheredDataReqNotFound = "no ciphered data found on request body"
)

var (
	// PayloadCryptoStrict will decide level of strictness of payload encryption mode.
	// `strict` will only accept/respond with ciphered data and time. Return error if ciphered data not found or un-decrypt-able.
	// `nonstrict` will accept/respond chipered data, time, and plain data. Return error when ciphered data exist and un-decrypt-able.
	PayloadCryptoStrict = condutils.Or(os.Getenv("PAYLOAD_CRYPTO_STRICT"), "").(string)

	// PayloadCryptoFeature set env var to `1` to enable payload crypto hooks on request/response.
	PayloadCryptoFeature = os.Getenv("PAYLOAD_CRYPTO")

	// PayloadCryptoKey to set key used in payload crypto.
	PayloadCryptoKey = condutils.Or(os.Getenv("PAYLOAD_CRYPTO_KEY"), "880d8e7e9b4b787aa50a3917b09fc0ec").(string)

	// PayloadCryptoMinAppVer set env var to `current app version code` to enable payload crypto hooks on configured version.
	PayloadCryptoMinAppVer = os.Getenv("PAYLOAD_CRYPTO_MIN_APP_VERSION")
)

// EncryptedResponse encrypted response for strict mode of encryption payload mode.
type EncryptedResponse struct {
	Ciphered  string `json:"xdata"`
	Timestamp int64  `json:"xtime"`
}

type decryptorFn func(cipherText, iv string) ([]byte, error)

func PeekAndPopBodyRequest(sLev string, bodyReq interface{}, decryptor decryptorFn) ([]byte, error) {
	// In not in `strict` nor `nonstrict` level, we don't care about ciphered data, return back original response.
	if !(sLev == "strict" || sLev == "nonstrict") {
		return json.Marshal(bodyReq)
	}

	rf := reflect.ValueOf(bodyReq)
	xData := rf.MapIndex(reflect.ValueOf("xdata"))
	xTime := rf.MapIndex(reflect.ValueOf("xtime"))

	// In `strict` level, we expect ciphered data found, return error when not found.
	if sLev == "strict" && (condutils.IsEmpty(xData) || condutils.IsEmpty(xTime)) {
		return nil, errors.New(CipheredDataReqNotFound)
	}

	// In `nonstrict` level, if no chipered data found return back payload.
	if sLev == "nonstrict" && (condutils.IsEmpty(xData) || condutils.IsEmpty(xTime)) {
		return json.Marshal(bodyReq)
	}

	timestamp := int64(xTime.Interface().(float64))
	cipherText := xData.Interface().(string)

	return decryptor(cipherText, fmt.Sprintf("%v", timestamp))
}

func AttachResponseBody(src map[string]interface{}, respBody interface{}) map[string]interface{} {
	if condutils.IsStruct(respBody) {
		rf := reflect.ValueOf(reflect.Indirect(reflect.ValueOf(respBody)).Interface())
		fieldNum := rf.NumField()
		for i := 0; i < fieldNum; i++ {
			fieldName := rf.Type().Field(i).Name
			// fieldTag := strings.Split(reflect.TypeOf(respBody).Field(i).Tag.Get("json"), ",")[0]
			fieldTag := strings.Split(rf.Type().Field(i).Tag.Get("json"), ",")[0]
			src[condutils.Or(fieldTag, fieldName).(string)] = rf.Field(i).Interface()
		}
		return src
	}

	mapIter := reflect.ValueOf(respBody).MapRange()
	for {
		if !mapIter.Next() {
			return src
		}
		src[mapIter.Key().String()] = mapIter.Value().Interface()
	}
}

func respPayloadHook(sLevel, cipherText string, timestamp int64, payload interface{}) interface{} {
	if sLevel == "strict" {
		return EncryptedResponse{Ciphered: cipherText, Timestamp: timestamp}
	}
	mapPayload := map[string]interface{}{"xdata": cipherText, "xtime": timestamp}

	if payload == nil {
		return mapPayload
	}

	return AttachResponseBody(mapPayload, payload)
}

// PreRequestHooks hook function called when binding data to request body.
// This may be used on handler.
func PreRequestHooks(key string, payload interface{}) (res []byte, rErr error) {
	defer func() {
		if r := recover(); r != nil {
			rErr = fmt.Errorf("got %v on pre-request hooks", r)
		}
	}()

	if condutils.IsEmpty(payload) {
		return json.Marshal(payload)
	}

	decryptor := func(key string) decryptorFn {
		return func(ciphered, iv string) ([]byte, error) {
			return crypts.PayloadDecrypt(ciphered, key, iv)
		}
	}
	return PeekAndPopBodyRequest(PayloadCryptoStrict, payload, decryptor(key))
}

// PreResponseHooks hook function called when try to write response body.
// This function may manipulate response data structure before the data sent to client.
func PreResponseHooks(key string, payload interface{}) (res interface{}, rErr error) {
	defer func() {
		if r := recover(); r != nil {
			rErr = fmt.Errorf("got %v on pre-response hooks", r)
		}
	}()

	if payload == nil {
		return payload, nil
	}

	buff, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return nil, marshalErr
	}

	timestamp := time.Now().Local().Unix()

	cipherText, encErr := crypts.PayloadEncrypt(buff, key, fmt.Sprintf("%v", timestamp))

	if encErr != nil {
		return nil, encErr
	}

	return respPayloadHook(PayloadCryptoStrict, cipherText, timestamp, payload), nil
}
