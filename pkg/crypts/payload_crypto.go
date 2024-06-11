package crypts

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func DeriveKeySHA256(key string, len int) []byte {
	hashed := sha256.Sum256([]byte(key))
	if len < 1 {
		len = sha256.Size
	}
	return hashed[:len]
}

func DeriveKeyHexSHA256(key string, len int) []byte {
	hashed := sha256.Sum256([]byte(key))
	if len < 1 {
		len = sha256.Size
	}
	hexed := fmt.Sprintf("%x", hashed)
	return []byte(hexed[:len])
	// return hashed[:len]
}

func PayloadEncrypt(message []byte, key, iv string) (cipherText string, rErr error) {
	defer func() {
		if r := recover(); r != nil {
			rErr = fmt.Errorf("got %v when encrypting", r)
		}
	}()

	block, err := aes.NewCipher(DeriveKeyHexSHA256(key, sha256.Size))
	if err != nil {
		return "", err
	}

	buff, err := Pad([]byte(message), aes.BlockSize)
	blkChiper := cipher.NewCBCEncrypter(block, (DeriveKeyHexSHA256(iv, aes.BlockSize)))
	blkChiper.CryptBlocks(buff, buff)
	return base64.URLEncoding.EncodeToString(buff), err
}

func PayloadDecrypt(cipherText, key, iv string) (plainMessage []byte, rErr error) {
	defer func() {
		if r := recover(); r != nil {
			rErr = fmt.Errorf("got %v when decrypting", r)
		}
	}()

	buff, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(DeriveKeyHexSHA256(key, sha256.Size))
	if err != nil {
		return nil, err
	}

	blkChiper := cipher.NewCBCDecrypter(block, DeriveKeyHexSHA256(iv, aes.BlockSize))
	blkChiper.CryptBlocks(buff, buff)
	return Unpad(buff, aes.BlockSize)
}
