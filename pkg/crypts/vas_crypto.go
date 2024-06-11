package crypts

import (
	"crypto/rand"
	"fmt"
)

const saltedHexLen = 32
const saltedBytesLen = 128 / 8

func genRandHex(num int) (string, error) {
	b := make([]byte, num)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// VASEncrypt ...
func VASEncrypt(input, key, iv string) (string, error) {
	cipherText, encErr := EncCJSAES(input, key)
	if encErr != nil {
		return "", encErr
	}
	salt, randErr := genRandHex(saltedBytesLen)
	if randErr != nil {
		return "", randErr
	}

	return salt + iv + cipherText, nil
}

// VASDecrypt ...
func VASDecrypt(input, key, iv string) (string, error) {
	lenIV := len(iv)
	tIV := []byte(input)[saltedHexLen : saltedHexLen+lenIV]
	cipherText := []byte(input)[saltedHexLen+lenIV:]
	if string(tIV) != iv {
		return "", fmt.Errorf("expect %v, got %v", iv, string(tIV))
	}
	res, decErr := DecCJSAES(string(cipherText), key)
	if decErr != nil {
		return "", decErr
	}
	return string(res), nil
}
