package jwt

// JWTConfigurationNotFound ...
const JWTConfigurationNotFound = "No configuration key jwt found in secrets entry"

// JWTSecretNotFound ...
const JWTSecretNotFound = "No configuration key secrets found"

// FieldStringCheckErrorMsg ...
const FieldStringCheckErrorMsg = `expect field %s not empty, but got %#v on your bearer token, please ask your bearer token issuer`

// UnknownSubsType substype is empty string or we don't know the substype.
const UnknownSubsType = `unknown subscriber type, got %v from field %s or %s`

// HomeFiberSubsType fiber substype for special case handling.
const HomeFiberSubsType = "HOMEFIBER"

// HomeSatuSubsType fiber substype for special case handling.
const HomeSatuSubsType = "HOMESATU"

// JWTClaims standard claims for JWT
type JWTClaims struct {
	JWTID        string `json:"jti"`
	Issuer       string `json:"iss"`
	ClientID     string `json:"client_id"`
	Subject      string `json:"sub"`
	Audience     string `json:"aud"`
	SessionID    string `json:"sid"`
	IssuedAt     int64  `json:"iat,omitempty"`
	ExpiredAt    int64  `json:"exp"`
	MSISDN       string `json:"-"`
	SubsID       string `json:"-"`
	DeviceID     string `json:"-"`
	IsFirstLogin string `json:"-"`
}

// CIAMClaims struct to hold information of claims got from CIAM
type CIAMClaims struct {
	AtHash          string `json:"at_hash"`
	Subject         string `json:"sub"`
	AuditTrackingID string `json:"auditTrackingId"`
	Issuer          string `json:"iss"`
	TokenName       string `json:"tokenName"`
	SubscriberID    string `json:"subscriberID"`
	Nonce           string `json:"nonce"`
	Audience        string `json:"aud"`
	CHash           string `json:"c_hash"`
	Acr             string `json:"acr"`
	ForgerockOpenID string `json:"org.forgerock.openidconnect.ops"`
	SHash           string `json:"s_hash"`
	Azp             string `json:"azp"`
	AuthTime        int64  `json:"auth_time"`
	Name            string `json:"name"`
	Realm           string `json:"realm"`
	MSISDN          string `json:"msisdn"`
	ExpiredAt       int64  `json:"exp"`
	TokenType       string `json:"tokenType"`
	FamilyName      string `json:"family_name"`
	IssuedAt        int64  `json:"iat"`
	DeviceID        string `json:"deviceID"`
	SubsType        string `json:"subscription_type"`
	IsFirstLogin    string `json:"isFirstLogin"`

	// Since implementation Is Enterprise check on login.
	SubsStatus    string `json:"subscriberStatus"`
	CustType      string `json:"customerType"`
	AccountNumber string `json:"accountNumber"`
	PricePlanID   string `json:"priceplanID"`

	// Since HOMEFIBER Special Case Bearer Token Handling.
	SubsID         string `json:"subscriber-id"`
	SubscriberType string `json:"subscriber-type"`
	AccountID      string `json:"accountID"`
	CustomerID     string `json:"customerID"`
	Email          string `json:"email"`
}

// IJWT jwt interface to generate token
type IJWT interface {
	WithMapClaims(map[string]interface{}) (string, error)
	WithStructClaims(JWTClaims) (string, error)
	Duration() int
	NewClaims() JWTClaims
	Serialize(compactedJWT string) (JWTClaims, error)
	Bind(compactedJWT string, customClaims interface{}) error
	GetRaw(compactedJWT string) ([]byte, error)
	Assert(JWTClaims) error
}
