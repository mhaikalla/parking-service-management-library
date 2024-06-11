package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
}

func (jts *JWTTestSuite) TestPanicIntialization() {
	jts.Panics(func() { NewJWT(map[string]map[string]interface{}{}) }, "panicked when empty configuration provided")

	jts.Panics(func() {
		NewJWT(map[string]map[string]interface{}{
			"jwt": {
				"signing_method": "HS256",
				"enable":         true,
				"duration":       60,
			},
		})
	}, "panicked when no secrets configuration provided")

	jts.Panics(func() {
		NewJWT(map[string]map[string]interface{}{
			"jwt": {
				"signing_method": "HS256",
				"enable":         true,
				"duration":       60,
			},
			"secrets": {},
		})
	}, "panicked when empty secrets configuration provided")

	jts.Panics(func() {
		NewJWT(map[string]map[string]interface{}{
			"jwt": {
				"signing_method": "HS256",
				"enable":         true,
				"duration":       "60",
			},
		})
	}, "panicked when try to coerce interface that don't have the same type")
}

func (jts *JWTTestSuite) TestGenerateJWS() {

	jws := NewJWT(map[string]map[string]interface{}{
		"jwt": {
			"signing_method": "HS256",
			"enable":         true,
			"duration":       60,
		},
		"secrets": {
			"jwt": map[string]interface{}{
				"key": "NRKqQdQ9pE0NLDPeUshePA==",
			},
		},
	})

	mapClaims := map[string]interface{}{
		"jti":       "62334859595823",
		"client_id": "aklsjdlkajslkdjlkja",
		"aud":       "kasjdlkjllala",
		"iss":       "Middleware",
		"sub":       "AutoLogin",
		"iat":       time.Now().Unix(),
		"ip":        "127.0.0.1",
		"sid":       "6/M4j43UPodQg0rQ1Dau/jo2LcWQxo6gIiS1XcNXFjM=",
	}

	jts.Equal(jws.Duration(), 60, "duration must return set on config")
	token, err := jws.WithMapClaims(mapClaims)
	jts.NoError(err, "must no error")
	jts.NotEmpty(token, "token must exists")

	s, err := jws.Serialize(token)

	jts.NoError(err, "no error when serializing")

	jts.NoError(jws.Assert(s), "assertion must not error")

	jts.Equal("62334859595823", s.JWTID, "jit must expected")
	jts.Equal("Middleware", s.Issuer, "iss must expected")

	badClaims := map[string]interface{}{
		"jti": math.Inf(-1),
	}

	_, err = jws.WithMapClaims(badClaims)
	jts.Error(err, "error serialize claims")

	badKey := NewJWT(map[string]map[string]interface{}{
		"jwt": {
			"signing_method": "HS256",
			"enable":         true,
			"duration":       60,
		},
		"secrets": {
			"jwt": map[string]interface{}{
				"key": "",
			},
		},
	})

	_, err = badKey.WithMapClaims(mapClaims)
	jts.Error(err, "must error because key was emptied")

	const badToken = `eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjE1OTEwODUxMjIsImlhdCI6MTU5MTA4MTUyMiwiaXAiOiIxMjcuMC4wLjEiLCJpc3MiOiJNaWRkbGV3YXJlIiwianRpIjoiNjIzMzQ4NTk1OTU4MjMiLCJyZXFpZCI6IjYvTTRqNDNVUG9kUWcwclExRGF1L2pvMkxjV1F4bzZnSWlTMVhjTlhGak09Iiwic3ViIjoiQXV0b0xvZ2luIn0.Q61tFZwv6dnESHV_ZQWfjw1Uw1dTueZn8bCVGTXNA`

	_, err = jws.Serialize(badToken)
	jts.Error(err, "must error, bad signature")

	const notJSONClaims = "eyJhbGciOiJIUzI1NiJ9.PHA-bG9yZW0gaXBzdW08L3A-.wn4U0ZhpIsDXcsIB-6LhIIavq0kgpDBRGAd0zIAS3CI"

	emptyClaims, err := jws.Serialize(notJSONClaims)
	jts.Error(err, "must error, not json claims")

	jts.Error(jws.Assert(emptyClaims), "assertion must error")

	claims := jws.NewClaims()

	claims.ClientID = "62334859595823"
	claims.Audience = "aklsjdlkajslkdjlkja"
	claims.Issuer = "Middleware"
	claims.Subject = "AutoLogin"
	claims.SessionID = "6/M4j43UPodQg0rQ1Dau/jo2LcWQxo6gIiS1XcNXFjM="

	resT, err := jws.WithStructClaims(claims)

	jts.NoError(err, "generate token must no error")

	_, err = jws.Serialize(resT)

	jts.NoError(err, "serializing must no error")

	buffRaw, rawReadErr := jws.GetRaw(token)
	jts.NoError(rawReadErr)
	jts.NotEmpty(buffRaw)
	buffRaw2, rawReadErr2 := jws.GetRaw("foobar")
	jts.Error(rawReadErr2)
	jts.Empty(buffRaw2)

	var binder map[string]interface{}

	jts.Error(jws.Bind("foobar", &binder))
	jts.NoError(jws.Bind(token, &binder))
	jts.NotEmpty(binder)
}

func (jts *JWTTestSuite) TestGenerateJWE() {

	jwe := NewJWT(map[string]map[string]interface{}{
		"jwt": {
			"encryption_method":  "A128CBC-HS256",
			"compression_method": "none",
			"signing_method":     "HS256",
			"enable":             true,
			"duration":           60,
		},
		"secrets": {
			"jwt": map[string]interface{}{
				"key":  "NRKqQdQ9pE0NLDPeUshePA==",
				"salt": "tUykjBxMk1w593Ng3ercZQ==",
			},
		},
	})

	mapClaims := map[string]interface{}{
		"jti":       "62334859595823",
		"client_id": "aklsjdlkajslkdjlkja",
		"aud":       "kasjdlkjllala",
		"iss":       "Middleware",
		"sub":       "AutoLogin",
		"iat":       time.Now().Unix(),
		"ip":        "127.0.0.1",
		"sid":       "6/M4j43UPodQg0rQ1Dau/jo2LcWQxo6gIiS1XcNXFjM=",
	}

	jts.Equal(jwe.Duration(), 60, "duration must return set on config")
	token, err := jwe.WithMapClaims(mapClaims)
	jts.NoError(err, "must no error")
	jts.NotEmpty(token, "token must exists")

	s, err := jwe.Serialize(token)

	jts.NoError(jwe.Assert(s), "assertion must not error")

	jts.NoError(err, "no error when serializing")

	jts.Equal("62334859595823", s.JWTID, "jit must expected")
	jts.Equal("Middleware", s.Issuer, "iss must expected")

	const badToken = `eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjE1OTEwODUxMjIsImlhdCI6MTU5MTA4MTUyMiwiaXAiOiIxMjcuMC4wLjEiLCJpc3MiOiJNaWRkbGV3YXJlIiwianRpIjoiNjIzMzQ4NTk1OTU4MjMiLCJyZXFpZCI6IjYvTTRqNDNVUG9kUWcwclExRGF1L2pvMkxjV1F4bzZnSWlTMVhjTlhGak09Iiwic3ViIjoiQXV0b0xvZ2luIn0.Q61tFZwv6dnESHV_ZQWfjw1Uw1dTueZn8bCVGTXNA`

	_, err = jwe.Serialize(badToken)
	jts.Error(err, "must error, bad signature")

	const notJSONClaims = `eyJhbGciOiJBMTI4S1ciLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0.ERMuQQE-YffQo__3JUhTTHHhLeqGgL99v-yKjdpSqMqDUebTFl-Hvi6dkLWPWEOs3YoNXcs-S0Y6EhHKkTWxh26MqnHz0T4C.GmkSyUg-lDIyYgtRILov1A.s3SjG-Kl-krBmSH-n55-8u8KSgkwWTy5oxhk05RU-F8.FP5GDrvwjb-UfHxkIIBRDA`

	emptyClaims, err := jwe.Serialize(notJSONClaims)
	jts.Error(err, "must error, not json claims")

	jts.Error(jwe.Assert(emptyClaims), "assertion must error")

	claims := jwe.NewClaims()

	claims.ClientID = "62334859595823"
	claims.Audience = "aklsjdlkajslkdjlkja"
	claims.Issuer = "Middleware"
	claims.Subject = "AutoLogin"
	claims.SessionID = "6/M4j43UPodQg0rQ1Dau/jo2LcWQxo6gIiS1XcNXFjM="

	resT, err := jwe.WithStructClaims(claims)

	jts.NoError(err, "generate token must no error")

	_, err = jwe.Serialize(resT)

	jts.NoError(err, "serializing must no error")

}

func (jts *JWTTestSuite) TestSigningMethodSelection() {
	jts.Equal("HS384", SigningMethodSelection("HS384").String(), "must as expected")
	jts.Equal("HS512", SigningMethodSelection("HS512").String(), "must as expected")
	jts.Equal("ES256", SigningMethodSelection("ES256").String(), "must as expected")
	jts.Equal("ES384", SigningMethodSelection("ES384").String(), "must as expected")
	jts.Equal("ES512", SigningMethodSelection("ES512").String(), "must as expected")
	jts.Equal("HS256", SigningMethodSelection("").String(), "must as expected")

}

func (jts *JWTTestSuite) TestDeriveKey() {
	const key = "foobar"
	const salt = "foobaz"

	jts.Equal(32, len(deriveKey(key, salt, 32)), "length of derivekey must expected")
	jts.Equal(64, len(deriveKey(key, salt, 64)), "length of derivekey must expected")

	res1 := deriveKey(key, salt, 64)
	res2 := deriveKey(key, salt, 64)

	jts.Equal(res1, res2, "same args must have same result")

}

func (jts *JWTTestSuite) TestEncryptionMethodPairs() {
	encMethod, keyEnc, key := EncryptionMethodPairs("A192CBC-HS384", "foobar", "foobaz")
	jts.Equal("A192CBC-HS384", encMethod.String(), "encryption method must expected")
	jts.Equal("A192KW", keyEnc.String(), "key encryption method must expected")
	jts.Equal(48, len(key.([]byte)), "length of key must expected")

	encMethod, keyEnc, key = EncryptionMethodPairs("A256CBC-HS512", "foobar", "foobaz")
	jts.Equal("A256CBC-HS512", encMethod.String(), "encryption method must expected")
	jts.Equal("A256KW", keyEnc.String(), "key encryption method must expected")
	jts.Equal(64, len(key.([]byte)), "length of key must expected")

	encMethod, keyEnc, key = EncryptionMethodPairs("A128GCM", "foobar", "foobaz")
	jts.Equal("A128GCM", encMethod.String(), "encryption method must expected")
	jts.Equal("A128KW", keyEnc.String(), "key encryption method must expected")
	jts.Equal(32, len(key.([]byte)), "length of key must expected")

	encMethod, keyEnc, key = EncryptionMethodPairs("A192GCM", "foobar", "foobaz")
	jts.Equal("A192GCM", encMethod.String(), "encryption method must expected")
	jts.Equal("A192KW", keyEnc.String(), "key encryption method must expected")
	jts.Equal(48, len(key.([]byte)), "length of key must expected")

	encMethod, keyEnc, key = EncryptionMethodPairs("A256GCM", "foobar", "foobaz")
	jts.Equal("A256GCM", encMethod.String(), "encryption method must expected")
	jts.Equal("A256KW", keyEnc.String(), "key encryption method must expected")
	jts.Equal(64, len(key.([]byte)), "length of key must expected")

	encMethod, keyEnc, key = EncryptionMethodPairs("", "foobar", "foobaz")
	jts.Equal("A128CBC-HS256", encMethod.String(), "encryption method must expected")
	jts.Equal("A128KW", keyEnc.String(), "key encryption method must expected")
	jts.Equal(32, len(key.([]byte)), "length of key must expected")
}

func (jts *JWTTestSuite) TestCompressionMethodSelection() {
	jts.Equal("DEF", CompressionMethodSelection("deflate").String(), "method selection must expected")
	jts.Equal("", CompressionMethodSelection("none").String(), "method selection must expected")
	jts.Equal("", CompressionMethodSelection("").String(), "method selection must expected")
}

func (jts *JWTTestSuite) TestCreateNewClaims() {
	j := NewJWT(map[string]map[string]interface{}{
		"jwt": {
			"signing_method": "HS256",
			"enable":         true,
			"duration":       60,
		},
		"secrets": {
			"jwt": map[string]interface{}{
				"key": "NRKqQdQ9pE0NLDPeUshePA==",
			},
		},
	})
	jts.NotEmpty(j.NewClaims().JWTID, "jit must predefined")
	jts.NotZero(j.NewClaims().IssuedAt, "iat must predefined")
}

func (jts *JWTTestSuite) TestDefaultAssertFn() {
	expiredToken := JWTClaims{ExpiredAt: time.Now().Unix()}
	jts.Error(defaultAssertFn(expiredToken), "must error")

	noClientID := JWTClaims{ExpiredAt: time.Now().Add(time.Second * 2).Unix()}
	jts.Error(defaultAssertFn(noClientID), "must error")

	noSubject := JWTClaims{ExpiredAt: time.Now().Add(time.Second * 2).Unix(), ClientID: "foobar"}
	jts.Error(defaultAssertFn(noSubject), "must error")

	noAudience := JWTClaims{ExpiredAt: time.Now().Add(time.Second * 2).Unix(), ClientID: "foobar", Subject: "foobar"}
	jts.Error(defaultAssertFn(noAudience), "must error")

	noID := JWTClaims{ExpiredAt: time.Now().Add(time.Second * 2).Unix(), ClientID: "foobar", Subject: "foobar", Audience: "bar"}
	jts.Error(defaultAssertFn(noID), "must error")

	noSessionID := JWTClaims{ExpiredAt: time.Now().Add(time.Second * 2).Unix(), ClientID: "foobar", Subject: "foobar", Audience: "bar", JWTID: "1"}
	jts.Error(defaultAssertFn(noSessionID), "must error")

	noIssuedAt := JWTClaims{ExpiredAt: time.Now().Add(time.Second * 2).Unix(), ClientID: "foobar", Subject: "foobar", Audience: "bar", JWTID: "1", SessionID: "0"}
	jts.Error(defaultAssertFn(noIssuedAt), "must error")

	complete := JWTClaims{ExpiredAt: time.Now().Add(time.Second * 2).Unix(), ClientID: "foobar", Subject: "foobar", Audience: "bar", JWTID: "1", IssuedAt: 1222332}
	jts.Error(defaultAssertFn(complete), "must error")
}

func (jts *JWTTestSuite) TestRSAJWE() {

	privateKey, privateKeyErr := rsa.GenerateKey(rand.Reader, 2048)
	jts.NoError(privateKeyErr)
	marshaled, mErr := x509.MarshalPKCS8PrivateKey(privateKey)
	jts.NoError(mErr)

	pkBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: marshaled,
	}

	pemBytes := pem.EncodeToMemory(pkBlock)

	configJWE := map[string]map[string]interface{}{
		"jwt": {
			"encryption_method":  "A128CBC-HS256",
			"key_algo":           "RSA-OAEP-256",
			"compression_method": "none",
			"enable":             true,
			"duration":           60,
		},
		"secrets": {
			"jwt": map[string]interface{}{
				"key": string(pemBytes),
			},
		},
	}

	rsaJWE := NewJWT(configJWE)

	claims := rsaJWE.NewClaims()
	claims.Audience = "foo"
	claims.ClientID = "foo1234"
	claims.Issuer = "some"
	claims.JWTID = uuid.New().String()
	claims.SessionID = uuid.New().String()
	claims.Subject = "test"

	jweToken, jweTokenErr := rsaJWE.WithStructClaims(claims)
	jts.NoError(jweTokenErr)

	marshaledToken, marshaledErr := rsaJWE.Serialize(jweToken)
	jts.NoError(marshaledErr)

	jts.Equal(claims.Audience, marshaledToken.Audience)
	jts.Equal(claims.ClientID, marshaledToken.ClientID)

}

func (jts *JWTTestSuite) TestRSAJWS() {

	privateKey, privateKeyErr := rsa.GenerateKey(rand.Reader, 2048)
	jts.NoError(privateKeyErr)
	marshaled, mErr := x509.MarshalPKCS8PrivateKey(privateKey)
	jts.NoError(mErr)

	pkBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: marshaled,
	}

	pemBytes := pem.EncodeToMemory(pkBlock)

	configJWS := map[string]map[string]interface{}{
		"jwt": {
			"signing_method": "RS256",
			"enable":         true,
			"duration":       60,
		},
		"secrets": {
			"jwt": map[string]interface{}{
				"key": string(pemBytes),
			},
		},
	}

	rsaJWE := NewJWT(configJWS)

	claims := rsaJWE.NewClaims()
	claims.Audience = "foo"
	claims.ClientID = "foo1234"
	claims.Issuer = "some"
	claims.JWTID = uuid.New().String()
	claims.SessionID = uuid.New().String()
	claims.Subject = "test"

	jweToken, jweTokenErr := rsaJWE.WithStructClaims(claims)
	jts.NoError(jweTokenErr)

	marshaledToken, marshaledErr := rsaJWE.Serialize(jweToken)
	jts.NoError(marshaledErr)

	jts.Equal(claims.Audience, marshaledToken.Audience)
	jts.Equal(claims.ClientID, marshaledToken.ClientID)
}

func TestJWTTestSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}
