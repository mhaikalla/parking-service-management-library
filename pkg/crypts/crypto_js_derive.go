package crypts

import (
	"crypto/aes"
	"crypto/cipher"

	/* #nosec G501 */
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

const (
	EVPIterCryptoJS       = 1
	OpenSSLCompatibleHead = `Salted__`
	SaltLen               = 8
	EVPLenKeys            = 16
	AES128KeyLen          = 32
)

// EVPBytesToKey implements the Openssl EVP_BytesToKey logic.
func EVPBytesToKey(keyLen, ivLen, nIter int, salt, password []byte, hasher hash.Hash) (key, iv []byte, err error) {

	hashLen := hasher.BlockSize()
	expectLen := (keyLen + ivLen + hashLen - 1) / hashLen * hashLen
	lastHash := []byte{}
	joined := []byte{}

	for ; len(joined) < expectLen; hasher.Reset() {
		_, wErr := hasher.Write(append(lastHash, append(password, salt...)...))
		if wErr != nil {
			return nil, nil, wErr
		}
		lastHash = hasher.Sum(nil)

		for i := 1; i < nIter; i++ {
			hasher.Reset()
			lastHash = hasher.Sum(lastHash)
			hasher.Reset()
		}

		joined = append(joined, lastHash...)
	}

	return joined[:keyLen], joined[keyLen : keyLen+ivLen], nil
}

// DecCJSAES ...
func DecCJSAES(cipherText, key string) (string, error) {
	keyLen := len(key)
	buff, base64Err := base64.StdEncoding.DecodeString(cipherText)
	if base64Err != nil {
		return "", base64Err
	}

	saltData := buff[SaltLen:aes.BlockSize]

	/* #nosec G401 */
	keyD, ivD, errD := EVPBytesToKey(keyLen, aes.BlockSize, EVPIterCryptoJS, saltData, []byte(key), sha256.New())

	if errD != nil {
		return "", errD
	}

	cText := buff[aes.BlockSize:]

	block, err := aes.NewCipher(keyD)
	if err != nil {
		return "", err
	}
	blkChiper := cipher.NewCBCDecrypter(block, ivD)
	blkChiper.CryptBlocks(cText, cText)
	plain, err := Unpad(cText, aes.BlockSize)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// EncCJSAES ...
func EncCJSAES(plainText, key string) (string, error) {

	keyLen := len(key)
	head := []byte(OpenSSLCompatibleHead)
	saltData := make([]byte, SaltLen)

	if _, randErr := rand.Read(saltData); randErr != nil {
		return "", randErr
	}

	/* #nosec G401 */
	keyD, ivD, errD := EVPBytesToKey(keyLen, aes.BlockSize, EVPIterCryptoJS, saltData, []byte(key), sha256.New())

	if errD != nil {
		return "", errD
	}

	block, cipherErr := aes.NewCipher(keyD)
	if cipherErr != nil {
		return "", cipherErr
	}

	buff, padErr := Pad([]byte(plainText), aes.BlockSize)
	if padErr != nil {
		return "", padErr
	}

	blkChiper := cipher.NewCBCEncrypter(block, ivD)
	blkChiper.CryptBlocks(buff, buff)

	head = append(head, saltData...)
	head = append(head, buff...)

	return base64.StdEncoding.EncodeToString(head), nil
}
