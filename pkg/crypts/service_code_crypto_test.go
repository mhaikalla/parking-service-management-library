package crypts

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type textCTX struct {
	msisdn   string
	subsType string
	subsID   string
}

func (tc textCTX) GetMSISDN() string {
	return tc.msisdn
}

func (tc textCTX) GetSubsID() string {
	return tc.subsID
}

func (tc textCTX) GetSubsType() string {
	return tc.subsType
}

func TestServiceCodeEncryption(t *testing.T) {
	ctx := &textCTX{msisdn: "123456789009876", subsID: "098767", subsType: "FOOBAR"}

	scc := NewServiceCodeCrypto(ctx)
	res, resErr := scc.Pack("1234", "", "")
	assert.NoError(t, resErr)

	res2, res2Err := scc.UnPack(res)
	assert.NoError(t, res2Err)
	assert.NotEmpty(t, res2)
	assert.Equal(t, "1234", res2.GetServiceCode())

	res3, res3Err := scc.Pack("908876", "99887", "90876789")
	assert.NoError(t, res3Err)

	res4, res4Err := scc.UnPack(res3)
	assert.NoError(t, res4Err)

	expirySess = -1
	defer func() {
		expirySess = 60
	}()

	assert.Equal(t, "908876", res4.GetServiceCode())
	assert.Equal(t, "99887", res4.GetOfferCode())
	assert.Equal(t, "90876789", res4.GetChannel())
	assert.False(t, res4.GetTimestamp().IsZero())
	assert.True(t, res4.IsExpired())
	assert.True(t, res4.IsOffer())

	validator := createExpiryValidator(10)
	assert.NotEmpty(t, validator)
	assert.True(t, validator(time.Now().Local()))
	assert.False(t, validator(time.Date(2020, 10, 02, 0, 0, 0, 0, time.Local)))
}

func Test_getCallerFunc(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test 1", "testing.tRunner"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCallerFunc(); got != tt.want {
				t.Errorf("getCallerFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_getCallerFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if getCallerFunc() != "testing.(*B).runN" {
			b.Fatal()
		}
	}
}

func Test_saltKey(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"test 0 args", args{keys: []string{}}, []byte("")},
		{"test 1 args", args{keys: []string{"foo"}}, []byte("")},
		{
			"test 2 args",
			args{keys: []string{"foo", "bar"}},
			[]byte{0x8b, 0x3f, 0x3e, 0x54, 0x5f, 0x80, 0x73, 0x96, 0x5, 0x61, 0x29, 0x6c, 0x20, 0xc7, 0x47, 0x58},
		},
		{
			"test 3 args",
			args{keys: []string{"foo", "bar", "baz"}},
			[]byte{0x90, 0x9f, 0xfe, 0xca, 0xef, 0x24, 0x5a, 0xa8, 0xfe, 0x93, 0x9c, 0x42, 0xf3, 0x8, 0xc, 0x83},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if fmt.Sprintf("%v", r) != "at least 2 arg must be supplied" {
						t.Fatal(r)
					}
				}
			}()
			if got := saltKey(tt.args.keys...); string(got) != string(tt.want) {
				t.Errorf("saltKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
