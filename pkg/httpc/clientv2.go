package httpc

import (
	"context"
	"net/http"
	"strings"

	"github.com/mhaikalla/parking-service-management-library/pkg/condutils"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
	"github.com/mhaikalla/parking-service-management-library/pkg/logs"
)

const (
	// RequestBodyError ...
	RequestBodyError = "CANNOT_PROCESSING_REQUEST_BODY"
	// RequestClientError ...
	RequestClientError = "CANNOT_PROCESSING_REQUEST_CLIENT"
	// ResponseUnserializeError ...
	ResponseUnserializeError = "CANNOT_PROCESSING_RESPONSE"
	// ResponseError ...
	ResponseError = "CANNOT_PROCESSING_REQUEST"
	// GenericError ..
	GenericError = "Code=%v; Status=%v; Message=%s; ReadError=%v"
)

// UpstreamsRequest ...
type UpstreamsRequest struct {
	URL                string
	Method             string
	Headers            map[string]string
	Bearer             string
	BearerFn           func() (string, error)
	Client             *http.Client
	BodyPayload        interface{}
	BodySuccess        interface{}
	BodyFailed         interface{}
	SuccessCodes       []int
	InsecureSkipVerify bool // set InsecureSkipVerify true for next upstream call implementation
	success            bool
	statusCode         int
	requestError       *errs.Errs
	defaultLogger      logs.ILog
	responseHeaders    http.Header
}

// IRequest interface to call to outside endpoints.
type IRequest interface {
	RequestWithContext(ctx context.Context) *UpstreamsRequest
	Request() *UpstreamsRequest
	GetError() *errs.Errs
	GetCode() int
	IsSuccess() bool
	GetHeaders() http.Header
}

// Request try requesting using defined UpstreamsRequest object
func (ur *UpstreamsRequest) Request() (returnedCtx *UpstreamsRequest) {
	return ur.RequestWithContext(context.TODO())
}

// RequestWithContext same like Request but using context provided
func (ur *UpstreamsRequest) RequestWithContext(ctx context.Context) *UpstreamsRequest {
	child, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	defer func() {
		condutils.When(ur.requestError != nil, ur.defaultLogger.Upsert, "error", ur.requestError)
		ur.defaultLogger.Upsert(logs.LogType, "OUTSIDE")
		ur.defaultLogger.Update()

		condutils.When(ur.requestError != nil, ur.defaultLogger.Error, "HTTP REQUEST")
		condutils.When(ur.requestError == nil, ur.defaultLogger.Info, "HTTP REQUEST")
	}()

	baggage := contexts.GetBaggage(child)
	ur.defaultLogger = baggage.Logger

	execute(child, ur)
	return ur
}

// GetError get error saved from previous operation
func (ur *UpstreamsRequest) GetError() *errs.Errs {
	return ur.requestError
}

// GetCode get status code of response
func (ur *UpstreamsRequest) GetCode() int {
	return ur.statusCode
}

// IsSuccess return true if no error occured, omit response status code
func (ur *UpstreamsRequest) IsSuccess() bool {
	return ur.success
}

// GetHeaders get response headers if exists.
func (ur *UpstreamsRequest) GetHeaders() http.Header {
	return ur.responseHeaders
}

// maskingBearer masking Bearer Token to avoid security concern.
func maskingBearer(bearer string) string {
	const bearerPrefix = "Bearer "
	bearer = strings.ReplaceAll(bearer, bearerPrefix, "")
	lenBearer := len(bearer)
	if lenBearer > 5 {
		if lenBearer%2 == 1 {
			halfed := (lenBearer - 3) / 2
			return bearerPrefix + bearer[:halfed] + "***" + bearer[(halfed+3):]
		}
		halfed := (lenBearer - 4) / 2
		return bearerPrefix + bearer[:halfed] + "****" + bearer[(halfed+4):]

	}
	return bearerPrefix + bearer
}

func maskString(str string) string {
	len := len(str)
	if len > 5 {
		str = str[0:3] + "***" + str[len-3:len]
		return str
	}
	return "***"
}
