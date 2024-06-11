package crypts

import (
	"crypto/sha256"
	"testing"
)

func TestEVPBytesToKey(t *testing.T) {
	key, iv, genErr := EVPBytesToKey(16, 16, 10, []byte("foobar"), []byte("password"), sha256.New())

	if genErr != nil || len(key) < 1 || len(iv) < 1 {
		t.Fail()
	}
}

func TestEncDecCJSAES(t *testing.T) {
	key := "1234567890987654"
	msg := "Hello World"

	cipherText, encErr := EncCJSAES(msg, key)
	if encErr != nil || cipherText == "" {
		t.Error(encErr)
	}

	plainText, decErr := DecCJSAES(cipherText, key)
	if decErr != nil || plainText != msg {
		t.Error(decErr)
	}

	if _, encErr := EncCJSAES("foobar", "009989"); encErr == nil {
		t.Fail()
	}

	if _, decErr := DecCJSAES(cipherText+"appended", key); decErr == nil {
		t.Fail()
	}
}
