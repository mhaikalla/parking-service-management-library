package crypts

import (
	"crypto/aes"
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestDeriveKeySHA256(t *testing.T) {
	type args struct {
		key string
		len int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"1", args{key: "foobar", len: 0}, 32},
		{"2", args{key: "foobar", len: 16}, 16},
		{"3", args{key: "foobar", len: -1}, 32},
		{"4", args{key: "", len: sha256.Size}, sha256.Size},
		{"5", args{key: "boo", len: aes.BlockSize}, aes.BlockSize},
		{"6", args{key: "boo", len: aes.BlockSize}, 16},
		{"7", args{key: "boo", len: sha256.Size}, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeriveKeySHA256(tt.args.key, tt.args.len); len(got) != tt.want {
				t.Errorf("DeriveKeySHA256() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestEncryptionDecryption(t *testing.T) {
	tmp1 := time.Now().Unix()
	m1 := []byte("foobar")
	iv1 := fmt.Sprintf("%v", tmp1)
	key1 := "988777898789"

	c1, ce1 := PayloadEncrypt(m1, key1, iv1)

	if ce1 != nil {
		t.FailNow()
	}

	d1, de1 := PayloadDecrypt(c1, key1, iv1)

	if de1 != nil {
		t.FailNow()
	}

	if !reflect.DeepEqual(d1, m1) {
		t.Error(string(d1))
	}

	if dek := DeriveKeySHA256("14021231241", sha256.Size); fmt.Sprintf("%x", dek) != "d4d5f187cda35b57d0e14c0fc199b9fdf754f099826a78af4d6ae308299bf128" {
		t.Error(fmt.Sprintf("%x", dek))
	}

	keyx := "14021231241"
	ivx := "1602088521655"
	cipherText := "x6gPi1t019hBPkDYTIFxzA=="

	if encc, _ := PayloadDecrypt(cipherText, keyx, ivx); string(encc) != "I am Iron Man." {
		t.Error(string(encc))
	}

	// ivv := "1604588855411"
	// ccp := "E1e8n8d605FcD2WNqxnFEirI6JqLNSMw3qA9x-6ohESHFQ_aot1j_HE-JFEfbjLmVw0kzdRaBN-3ShvfTi16lQ\u003d\u003d"

	// if enccccc, err := PayloadDecrypt(ccp, "880d8e7e9b4b787aa50a3917b09fc0ec", ivv); err != nil {
	// 	t.Error(err)
	// 	// t.Error(string(encc))
	// } else {
	// 	t.Error(string(enccccc))
	// }
}
