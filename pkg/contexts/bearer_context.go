package contexts

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"parking-service/pkg/jwt"
	"parking-service/pkg/logs"
	"parking-service/pkg/payloadhooks"

	"github.com/labstack/echo/v4"
)

var (
	serviceName   = os.Getenv("SERVICE_NAME")
	reBearerToken = regexp.MustCompile(`[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`)
)

// SideLoad struct containing Bearer Token.
type SideLoad struct {
	BearerData jwt.JWTClaims
	CiamClaims jwt.CIAMClaims
}

type RequestValidationData struct {
	XVersionApp   string
	UserAgent     string
	XSignature    string
	XHV           string
	RequestMethod string
	RequestPath   string
}

// BearerContext context passing from interceptor to handler.
type BearerContext struct {
	echo.Context
	SideLoad          SideLoad
	RequestID         string
	preReqHook        func(key string, reqPayload interface{}) ([]byte, error)
	preRespHook       func(key string, respPayload interface{}) (interface{}, error)
	preReqHookStrict  func(key string, reqPayload interface{}) ([]byte, error)
	preRespHookStrict func(key string, respPayload interface{}) (interface{}, error)
	hookEnable        bool
	logger            logs.ILog

	RequestValidationData
}

// JSON overide echo.Context#JSON method, will encrypt field Data
// when hookEnable is true and preRespHook is not nil.
func (bc BearerContext) JSON(code int, payload interface{}) error {
	bc.logger.Update()

	defer func() {
		bc.logger.Upsert("response", payload)
		bc.logger.Update()
		bc.logger.Debug("RESPONSE COMPLETED")
	}()

	// add validation to avoid error response header overwrite
	if bc.Response().Header().Get("status-code") == "" {
		bc.Response().Header().Set("status-code", "200")
	}

	if bc.hookEnable && bc.preRespHook != nil {
		key := getKeyFromEnv(bc)
		res, hookErr := bc.preRespHook(key, payload)
		if hookErr != nil {
			return hookErr
		}

		bc.logger.Upsert("encrypted_response", res)
		return bc.Context.JSON(code, res)
	}

	return bc.Context.JSON(code, payload)
}

// JSON overide echo.Context#JSON method, will encrypt field Data
// when hookEnable is true and preRespHook is not nil.
func (bc BearerContext) JSONEncrypted(code int, payload interface{}) error {
	bc.logger.Update()

	defer func() {
		bc.logger.Upsert("response", payload)
		bc.logger.Update()
		bc.logger.Debug("RESPONSE COMPLETED")
	}()

	// add validation to avoid error response header overwrite
	if bc.Response().Header().Get("status-code") == "" {
		bc.Response().Header().Set("status-code", "200")
	}

	key := getKeyFromEnv(bc)
	res, hookErr := bc.preRespHookStrict(key, payload)
	if hookErr != nil {
		return hookErr
	}

	bc.logger.Upsert("encrypted_response", res)
	return bc.Context.JSON(code, res)
}

// Load like echo.Context#Bind method, will try to decrypt Data field
// in request body payload when hookEnable is true and preReqHook is not nil.
func (bc BearerContext) Load(structure interface{}) error {
	if bc.hookEnable && bc.preReqHook != nil {
		var reqPayload interface{}
		if err := bc.Context.Bind(&reqPayload); err != nil {
			return err
		}

		key := getKeyFromEnv(bc)

		buff, hookErr := bc.preReqHook(key, reqPayload)
		if hookErr != nil {
			return hookErr
		}
		return json.Unmarshal(buff, &structure)
	}
	return bc.Context.Bind(&structure)
}

// LoadEncrypted like Load method, but only for encrypted use.
func (bc BearerContext) LoadEncrypted(structure interface{}) error {
	var reqPayload interface{}
	if err := bc.Context.Bind(&reqPayload); err != nil {
		return err
	}

	key := getKeyFromEnv(bc)
	buff, hookErr := bc.preReqHookStrict(key, reqPayload)
	if hookErr != nil {
		return hookErr
	}
	return json.Unmarshal(buff, &structure)
}

// Consume like parking-service/pkg/interceptors.BearerContext#Load method,
// but just validating payload and trigger preRespHook activation.
func (bc BearerContext) Consume() error {
	if bc.hookEnable && bc.preReqHook != nil {
		var reqPayload interface{}
		if err := bc.Context.Bind(&reqPayload); err != nil {
			return err
		}
		if reqPayload == nil {
			return nil
		}
		key := getKeyFromEnv(bc)
		_, hookErr := bc.preReqHook(key, reqPayload)
		if hookErr != nil {
			return hookErr
		}
	}
	return nil
}

// GetMSISDN get MSISDN from JWT claim.
func (bc *BearerContext) GetMSISDN() string {
	return bc.SideLoad.BearerData.MSISDN
}

// GetSubsID get subscriber ID from JWT Claim.
func (bc *BearerContext) GetSubsID() string {
	return bc.SideLoad.BearerData.SubsID
}

// GetDeviceID get device ID from JWT Claim.
func (bc *BearerContext) GetDeviceID() string {
	return bc.SideLoad.BearerData.DeviceID
}

// GetSubsType get subscriber type like PREPAID, GO, POSTPAID,etc.
func (bc *BearerContext) GetSubsType() string {
	return bc.SideLoad.BearerData.Audience
}

// GetRequestID get request ID from request header X-Request-ID.
func (bc *BearerContext) GetRequestID() string {
	return bc.RequestID
}

// GetRequestContext get request context.
func (bc *BearerContext) GetRequestContext() context.Context {
	baggage := Baggage{
		PathOrigin:  bc.Request().RequestURI,
		DeviceID:    bc.GetDeviceID(),
		RequestID:   bc.RequestID,
		SubsID:      bc.GetSubsID(),
		Substype:    bc.GetSubsType(),
		Logger:      bc.logger.Child("CLIENT REQUEST"),
		BearerToken: reBearerToken.FindString(bc.Request().Header.Get(echo.HeaderAuthorization)),
	}
	contextWithValues := context.WithValue(bc.Request().Context(), requestContextBaggage, baggage)
	return contextWithValues
}

// Logger return log wrapper.
func (bc BearerContext) GetLogger() logs.ILog {
	if bc.logger != nil {
		return bc.logger
	}
	return logs.NewLogrus("NEW LOG HANDLER")
}

// GetPricePlan from Bearer Token.
func (bc BearerContext) GetPricePlan() string {
	return bc.SideLoad.CiamClaims.PricePlanID
}

// GetCustomerType from Bearer Token.
func (bc BearerContext) GetCustomerType() string {
	return bc.SideLoad.CiamClaims.CustType
}

// GetSubscriberStatus from Bearer Token.
func (bc BearerContext) GetSubscriberStatus() string {
	return bc.SideLoad.CiamClaims.SubsStatus
}

// GetAccountNumber from Bearer Token.
func (bc BearerContext) GetAccountNumber() string {
	return bc.SideLoad.CiamClaims.AccountNumber
}

// GetAuthTime from Bearer Token.
func (bc BearerContext) GetAuthTime() time.Time {
	return time.Unix(bc.SideLoad.CiamClaims.AuthTime, 0)
}

// GetAccountID get Account ID from CIAM Bearer Token ID, use by HOMEFIBER.
// Since HOMEFIBER Special Case Bearer Token Handling.
func (bc *BearerContext) GetAccountID() string {
	return bc.SideLoad.CiamClaims.AccountID
}

// GetCustomerID get Customer ID from CIAM Bearer Token ID, use by HOMEFIBER.
// Since HOMEFIBER Special Case Bearer Token Handling.
func (bc *BearerContext) GetCustomerID() string {
	if bc.SideLoad.CiamClaims.Audience == "HOMEFIBER" || bc.SideLoad.CiamClaims.Audience == "HOMESATU" {
		return bc.SideLoad.CiamClaims.AccountID
	}
	return bc.SideLoad.CiamClaims.CustomerID
}

// GetCustomerID get Email from CIAM Bearer Token ID, use by HOMEFIBER.
// Since HOMEFIBER Special Case Bearer Token Handling.
func (bc *BearerContext) GetEmail() string {
	return bc.SideLoad.CiamClaims.Email
}

// GetFirstlogin get firstlogin from CIAM Bearer Token ID, use by APPSFLYER.
func (bc *BearerContext) GetFirstlogin() string {
	return bc.SideLoad.CiamClaims.IsFirstLogin
}

// EnsureBearerContext ensure output is `BearerContext`, if not, wrap it on `BearerContext`.
// will ensure payload hook enable when env `PAYLOAD_CRYPTO` set to `1`
func EnsureBearerContext(maybeContext interface{}) BearerContext {
	payloadCryptoEnabled := (payloadhooks.PayloadCryptoFeature == "1")
	if bearerContext, bcOK := maybeContext.(BearerContext); bcOK {
		versionCode := extractVersionCode(bearerContext.Request().Header.Get("user-agent"))
		pcMinAppVer, _ := strconv.Atoi(payloadhooks.PayloadCryptoMinAppVer)
		payloadCryptoEnabled = versionCode > pcMinAppVer && (payloadhooks.PayloadCryptoFeature == "1")
		if bearerContext.logger == nil {
			bearerContext.logger = logs.NewLogrus(serviceName)
		}
		return ensurePayloadHooks(bearerContext, payloadCryptoEnabled)
	}

	if echoContext, ecOK := maybeContext.(echo.Context); ecOK {
		versionCode := extractVersionCode(echoContext.Request().Header.Get("user-agent"))
		pcMinAppVer, _ := strconv.Atoi(payloadhooks.PayloadCryptoMinAppVer)
		payloadCryptoEnabled = versionCode > pcMinAppVer && (payloadhooks.PayloadCryptoFeature == "1")
	}

	return ensurePayloadHooks(
		BearerContext{
			Context: maybeContext.(echo.Context),
			logger:  logs.NewLogrus(serviceName),
		},
		payloadCryptoEnabled,
	)
}

// getKeyFromContext get key for encryption based on existence of bearer token.
func getKeyFromEnv(ctx BearerContext) string {
	return payloadhooks.PayloadCryptoKey
}

// ensurePayloadHooks ensure payload hook enable.
func ensurePayloadHooks(bc BearerContext, enable bool) BearerContext {
	if !bc.hookEnable && enable {
		bc.preReqHook = payloadhooks.PreRequestHooks
		bc.preRespHook = payloadhooks.PreResponseHooks
		bc.hookEnable = true
	}
	bc.preReqHookStrict = payloadhooks.PreRequestHooks
	bc.preRespHookStrict = payloadhooks.PreResponseHooks
	return bc
}

// extractVersionCode get version code from user-agent request header, if invalid return 9999 to ensure encrypted payload is enabled
func extractVersionCode(userAgent string) int {
	splittedUA := strings.Split(userAgent, ";")
	if !strings.Contains(strings.ToLower(splittedUA[0]), "myxl") {
		return 9999
	}

	explodeVersion := strings.Split(splittedUA[0], "/")
	explodeVersionCode := strings.Split(explodeVersion[1], "(")
	if len(explodeVersionCode) < 2 {
		explodeVersionCode = []string{"", "0", ""}
	}

	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	extractedVersionCode := nonAlphanumericRegex.ReplaceAllString(explodeVersionCode[1], "")
	ver, err := strconv.Atoi(extractedVersionCode)
	if err != nil {
		return 9999
	}

	return ver
}
