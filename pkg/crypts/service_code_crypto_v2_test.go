package crypts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setServiceCodeCryptoOn(t *testing.T) {
	ServiceCodeCryptoEnabled = "1"
}

func setServiceCodeCryptoOff(t *testing.T) {
	ServiceCodeCryptoEnabled = ""
}

func TestServiceCodeEncryptionV2EncryptedEnv(t *testing.T) {
	setServiceCodeCryptoOn(t)
	ctx := &textCTX{msisdn: "123456789009876", subsID: "098767", subsType: "FOOBAR"}

	scc := NewServiceCodeCryptoV2(ctx)

	assert.NotPanics(t, func() { scc.Pack("1234") })
	res := scc.Pack("1234")
	assert.NotPanics(t, func() { scc.UnPack(res) })
	res2 := scc.UnPack(res)
	assert.NotEmpty(t, res2)
	assert.Equal(t, "1234", res2.GetServiceCode())
	assert.NotEmpty(t, res2.GetTimestamp())

	assert.NotPanics(t, func() { scc.Pack("908876", "99887", "CHANNEL") })
	res3 := scc.Pack("908876", "99887", "CHANNEL")
	assert.NotPanics(t, func() { scc.UnPack(res3) })
	res4 := scc.UnPack(res3)
	assert.NotEmpty(t, res4)
	assert.True(t, res4.IsOffer())
	assert.False(t, res4.IsLocationExist())
	assert.Equal(t, "908876", res4.GetServiceCode())
	assert.Equal(t, "99887", res4.GetOfferCode())
	assert.Equal(t, "CHANNEL", res4.GetChannel())

	assert.NotPanics(t, func() { scc.Pack("908876", "99887", "CHANNEL", "CHANNEL*A*B*C") })
	res5 := scc.Pack("908876", "99887", "CHANNEL", "CHANNEL*A*B*C")
	assert.NotPanics(t, func() { scc.UnPack(res5) })
	res6 := scc.UnPack(res5)
	assert.NotEmpty(t, res6)
	assert.True(t, res6.IsOffer())
	assert.True(t, res6.IsLocationExist())
	assert.Equal(t, "908876", res6.GetServiceCode())
	assert.Equal(t, "99887", res6.GetOfferCode())
	assert.Equal(t, "CHANNEL", res6.GetChannel())
	assert.Equal(t, "CHANNEL*A*B*C", res6.GetLocation())

	assert.NotPanics(t, func() { scc.Pack("908876|99887|CHANNEL|CHANNEL*A*B*C") })
	res7 := scc.Pack("908876|99887|CHANNEL|CHANNEL*A*B*C")
	assert.NotPanics(t, func() { scc.UnPack(res7) })
	res8 := scc.UnPack(res7)
	assert.NotEmpty(t, res8)
	assert.True(t, res8.IsOffer())
	assert.True(t, res8.IsLocationExist())
	assert.Equal(t, "908876", res8.GetServiceCode())
	assert.Equal(t, "99887", res8.GetOfferCode())
	assert.Equal(t, "CHANNEL", res8.GetChannel())
	assert.Equal(t, "CHANNEL*A*B*C", res8.GetLocation())

	assert.NotPanics(t, func() { scc.Pack("908876|99887|CHANNEL") })
	res9 := scc.Pack("908876|99887|CHANNEL")
	assert.NotPanics(t, func() { scc.UnPack(res9) })
	res10 := scc.UnPack(res9)
	assert.NotEmpty(t, res10)
	assert.True(t, res10.IsOffer())
	assert.False(t, res10.IsLocationExist())
	assert.Equal(t, "908876", res10.GetServiceCode())
	assert.Equal(t, "99887", res10.GetOfferCode())
	assert.Equal(t, "CHANNEL", res10.GetChannel())

	// Unencrypted Env
	setServiceCodeCryptoOff(t)
	scc1 := NewServiceCodeCryptoV2(ctx)

	assert.NotPanics(t, func() { scc1.Pack("1234") })
	res11 := scc1.Pack("1234")
	assert.NotPanics(t, func() { scc1.UnPack(res11) })
	res12 := scc1.UnPack(res11)
	assert.NotEmpty(t, res12)
	assert.Equal(t, "1234", res12.GetServiceCode())
	assert.NotEmpty(t, res12.GetTimestamp())

	assert.NotPanics(t, func() { scc1.Pack("908876", "99887", "CHANNEL") })
	res13 := scc1.Pack("908876", "99887", "CHANNEL")
	assert.NotPanics(t, func() { scc1.UnPack(res13) })
	res14 := scc1.UnPack(res13)
	assert.NotEmpty(t, res14)
	assert.True(t, res14.IsOffer())
	assert.False(t, res14.IsLocationExist())
	assert.Equal(t, "908876", res14.GetServiceCode())
	assert.Equal(t, "99887", res14.GetOfferCode())
	assert.Equal(t, "CHANNEL", res14.GetChannel())

	assert.NotPanics(t, func() { scc1.Pack("908876", "99887", "CHANNEL", "CHANNEL*A*B*C") })
	res15 := scc1.Pack("908876", "99887", "CHANNEL", "CHANNEL*A*B*C")
	assert.NotPanics(t, func() { scc1.UnPack(res15) })
	res16 := scc1.UnPack(res15)
	assert.NotEmpty(t, res16)
	assert.True(t, res16.IsOffer())
	assert.True(t, res16.IsLocationExist())
	assert.Equal(t, "908876", res16.GetServiceCode())
	assert.Equal(t, "99887", res16.GetOfferCode())
	assert.Equal(t, "CHANNEL", res16.GetChannel())
	assert.Equal(t, "CHANNEL*A*B*C", res16.GetLocation())

	assert.NotPanics(t, func() { scc1.Pack("908876|99887|CHANNEL|CHANNEL*A*B*C") })
	res17 := scc1.Pack("908876|99887|CHANNEL|CHANNEL*A*B*C")
	assert.NotPanics(t, func() { scc1.UnPack(res17) })
	res18 := scc1.UnPack(res17)
	assert.NotEmpty(t, res18)
	assert.True(t, res18.IsOffer())
	assert.True(t, res18.IsLocationExist())
	assert.Equal(t, "908876", res18.GetServiceCode())
	assert.Equal(t, "99887", res18.GetOfferCode())
	assert.Equal(t, "CHANNEL", res18.GetChannel())
	assert.Equal(t, "CHANNEL*A*B*C", res18.GetLocation())

	expirySess = -1
	defer func() {
		expirySess = 60
	}()

	assert.NotPanics(t, func() { scc1.Pack("908876|99887|CHANNEL") })
	res19 := scc1.Pack("908876|99887|CHANNEL")
	assert.NotPanics(t, func() { scc1.UnPack(res19) })
	res20 := scc1.UnPack(res19)
	assert.NotEmpty(t, res20)
	assert.True(t, res20.IsOffer())
	assert.False(t, res20.IsLocationExist())
	assert.Equal(t, "908876", res20.GetServiceCode())
	assert.Equal(t, "99887", res20.GetOfferCode())
	assert.Equal(t, "CHANNEL", res20.GetChannel())
	assert.True(t, res20.IsExpired())
	assert.True(t, res20.IsOffer())

	validator := createExpiryValidator(10)
	assert.NotEmpty(t, validator)
	assert.True(t, validator(time.Now().Local()))
	assert.False(t, validator(time.Date(2020, 10, 02, 0, 0, 0, 0, time.Local)))
}
