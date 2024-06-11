package httpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mhaikalla/parking-service-management-library/pkg/condutils"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

// parseBodyPayload parse body payload based object type, return reader, length of payload, and an error.
func parseBodyPayload(bodyPayload interface{}) (io.Reader, int, error) {
	switch t := bodyPayload.(type) {

	case nil:
		return nil, 0, nil

	case string:
		return strings.NewReader(t), len([]byte(t)), nil

	case []byte:
		return bytes.NewReader(t), len(t), nil

	case io.Reader:
		return t, -1, nil

	default:
		buff, err := json.Marshal(bodyPayload)
		if err != nil {
			return nil, 0, err
		}
		return bytes.NewReader(buff), len(buff), nil
	}
}

// attachAuthHeader if `bearerFn` not nil, evaluate it and on success attach result to request header.
// if `bearerFn` is nil, use `bearer` instead.
func attachAuthHeader(headers http.Header, bearerFn func() (string, error), bearer string) error {
	if bearerFn != nil {
		generated, genErr := bearerFn()
		if genErr != nil {
			return genErr
		}
		headers.Add("Authorization", "Bearer "+generated)
		return nil
	}

	condutils.When(bearer != "", headers.Add, "Authorization", "Bearer "+bearer)
	return nil
}

// bindJSONByteBuffer binding data structure with JSON byte array.
func bindJSONByteBuffer(buff []byte, structure interface{}) (bool, *errs.Errs) {
	// we try to unserialize body request using defined structure when content length above zero
	if len(buff) > 0 && structure != nil {
		if unmarshalErr := json.Unmarshal(buff, &structure); unmarshalErr != nil {
			return false, errs.NewErrContext().SetCode(errs.HTTPClientRequestErr).SetError(unmarshalErr)
		}
	}
	return true, nil
}

func createRequest(ur *UpstreamsRequest, ctx context.Context) (*http.Request, error) {
	// parse body payload to io.reader and check for content length
	payload, payloadLen, parseErr := parseBodyPayload(ur.BodyPayload)

	if parseErr != nil {
		return nil, errs.NewErrContext().SetCode(errs.HTTPClientRequestErr).SetError(parseErr)
	}

	// set InsecureSkipVerify true for next upstream call implementation
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: ur.InsecureSkipVerify}
	req, createErr := http.NewRequest(ur.Method, ur.URL, payload)

	if createErr != nil {
		return nil, errs.NewErrContext().SetCode(errs.HTTPClientRequestErr).SetError(createErr)
	}

	// we assume to send json body payload
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(payloadLen))

	// add headers key val to the request
	if ur.Headers != nil {
		for k, v := range ur.Headers {
			req.Header.Set(k, v)
		}
	}

	if err := attachAuthHeader(req.Header, ur.BearerFn, ur.Bearer); err != nil {
		return nil, errs.NewErrContext().SetCode(errs.HTTPClientRequestErr).SetError(err)
	}

	return req.WithContext(ctx), nil
}

// execute `UpstreamsRequest` using context provided.
func execute(ctx context.Context, ur *UpstreamsRequest) {
	req, createReqErr := createRequest(ur, ctx)
	cl := ur.defaultLogger.Child("UPSTREAM_REQUEST_LOG")

	defer func() {
		cl.Update()
		cl.Info()
	}()

	if createReqErr != nil {
		ur.requestError = errs.NewErrContext().SetCode(errs.HTTPClientRequestErr).SetError(createReqErr)
		return
	}

	ur.defaultLogger.Upsert("uri", fmt.Sprintf("%v %v", req.URL, req.Proto))
	cl.Upsert("uri", fmt.Sprintf("%v %v", req.URL, req.Proto))
	ur.defaultLogger.Upsert("request_method", req.Method)
	cl.Upsert("request_method", req.Method)
	ur.defaultLogger.Upsert("request_body", ur.BodyPayload)

	client := ur.Client

	startTime := time.Now()
	resp, respErr := client.Do(req)

	latency := time.Since(startTime).Milliseconds()
	cloned := req.Header.Clone()
	masked := maskingBearer(cloned.Get("Authorization"))
	maskedKey := maskString(cloned.Get("api_key"))
	cloned.Set("Authorization", masked)
	cloned.Set("api_key", maskedKey)
	ur.defaultLogger.Upsert("request_header", cloned)
	ur.defaultLogger.Upsert("request_start_time", startTime.Unix())
	ur.defaultLogger.Upsert("response_latency", latency)
	cl.Upsert("response_latency", latency)

	if respErr != nil {
		ur.requestError = errs.NewErrContext().SetCode(errs.HTTPClientRequestErr).SetError(respErr)
		return
	}

	defer resp.Body.Close()

	ur.defaultLogger.Upsert("response_code", resp.StatusCode)
	cl.Upsert("response_code", resp.StatusCode)

	var bodyToSerialize interface{} = nil
	// check if response status code match defined status code and set it's body to BodySuccess
	// short circuited on first match code
	for _, code := range ur.SuccessCodes {
		if resp.StatusCode == code {
			bodyToSerialize = &ur.BodySuccess
			ur.statusCode = resp.StatusCode
			ur.success = true
			break
		}
	}

	buff, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		ur.requestError = errs.NewErrContext().SetCode(errs.HTTPClientResponseErr).SetError(readErr).SetHttpCode(resp.StatusCode)
		return
	}

	defer ur.defaultLogger.Upsert("response_body", string(buff))

	if !ur.success {
		ur.statusCode = resp.StatusCode
		// if BodyFailed is nil, mean we will return an error and discard response if any.
		if ur.BodyFailed == nil {
			respErr := fmt.Errorf(GenericError, resp.StatusCode, resp.Status, takeSomeBuff(buff), readErr)
			ur.requestError = errs.NewErrContext().SetCode(errs.HTTPClientResponseErr).SetError(respErr).SetHttpCode(resp.StatusCode)
			return
		}

		bodyToSerialize = &ur.BodyFailed
	}
	ur.responseHeaders = resp.Header.Clone()
	ur.success, ur.requestError = bindJSONByteBuffer(buff, bodyToSerialize)
	if ur.requestError != nil {
		ur.requestError.SetHttpCode(resp.StatusCode)
	}
}

// takeSomeBuff take some byte array capped to 100 char.
func takeSomeBuff(buff []byte) string {
	if len(buff) > 100 {
		return string(buff[:100])
	}
	return string(buff)
}

// createHTTPClientRequestErr create client request error.
// This happen on client request creation.
func createHTTPClientRequestErr(orginalError error) *errs.Errs {
	return errs.NewErrContext().SetCode(errs.HTTPClientRequestErr).SetError(orginalError)
}

// createHTTPClientRequestBodyErr create client request body error.
// This happen on client body request creation.
func createHTTPClientRequestBodyErr(originalError error) *errs.Errs {
	return errs.NewErrContext().SetCode(errs.HTTPCLientRequestBodyErr).SetError(originalError)
}

// createHTTPClientResponseErr create client response error.
// This happen when request executed or response not handled.
func createHTTPClientResponseErr(originalError error) *errs.Errs {
	return errs.NewErrContext().SetCode(errs.HTTPClientResponseErr).SetError(originalError)
}
