package crypts

import "testing"

func Test_genRandHex(t *testing.T) {
	if res, genErr := genRandHex(16); genErr != nil || len(res) != 32 {
		t.Error(res)
	}

	if res, genErr := genRandHex(8); genErr != nil || len(res) != 16 {
		t.Error(res)
	}
}

func TestVASEncryptDecryptSuccess(t *testing.T) {
	key := "0123456789009876"
	iv := "foobar"
	msg := "hello world"

	cipherText, encErr := VASEncrypt(msg, key, iv)
	if encErr != nil || cipherText == "" {
		t.Fatal(encErr)
	}

	plainText, decErr := VASDecrypt(cipherText, key, iv)
	if decErr != nil || plainText != msg {
		t.Fatal(decErr)
	}

	if _, encErr := VASEncrypt("hello world", "not incremental 16 char", "foobar"); encErr == nil {
		t.Fail()
	}

	if _, decErr := VASDecrypt("not-chipered-text-and-must-be-error", "1234567890987654", "iv"); decErr == nil {
		t.Fail()
	}
}
