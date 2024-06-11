package errs

import (
	"fmt"
	"os"
	"strconv"
)

var (
	DefaultErrTitle   = os.Getenv("DEFAULT_ERROR_TITLE_ID")
	DefaultErrDesc    = os.Getenv("DEFAULT_ERROR_DESCRIPTION_ID")
	DefaultErrTitleEN = os.Getenv("DEFAULT_ERROR_TITLE_EN")
	DefaultErrDescEN  = os.Getenv("DEFAULT_ERROR_DESCRIPTION_EN")
)

// Errs implementation of `errs.Erss`.
type Errs struct {
	Code          string                 `json:"code"`
	Status        string                 `json:"status"`
	Message       string                 `json:"message"`
	HttpCode      string                 `json:"-"`
	OrigError     error                  `json:"-"`
	Location      string                 `json:"-"`
	ErrorCodesMap map[string]interface{} `json:"-"`
	Lang          string                 `json:"-"`
}

// IErrs interfacing `Error` interface with additional data.
type IErrs interface {
	SetCode(int) *Errs
	SetHttpCode(int) *Errs
	SetMessage(string) *Errs
	SetError(error) *Errs
	SetLocation() *Errs
	Log(func(...interface{})) *Errs
	GetData() *Errs
	Error() string
}

// Error satisfy `Error` interface.
func (e *Errs) Error() string {
	return fmt.Sprintf(`code=%s;status=%s;message=%s;error=%v;location=%s`, e.Code, e.Status, e.Message, e.OrigError, e.Location)
}

// SetCode set the error code approriate to the current context.
func (e *Errs) SetCode(code int) *Errs {
	e.Code = strconv.Itoa(code)
	return e
}

// SetHttpCode set the error code approriate to the current context.
func (e *Errs) SetHttpCode(code int) *Errs {
	e.HttpCode = strconv.Itoa(code)
	return e
}

// SetMessage set message error apporiate to the current context.
func (e *Errs) SetMessage(message string) *Errs {
	e.Message = message
	return e
}

// SetError set original error.
func (e *Errs) SetError(err error) *Errs {
	e.OrigError = err
	if err != nil {
		e.Message = err.Error()
	}
	return e
}

// SetLocation set location when this error generated.
func (e *Errs) SetLocation() *Errs {
	e.Location = getCallerFunc()
	return e
}

// Log do logging of this error.
func (e *Errs) Log(logger func(...interface{})) *Errs {
	logger(e)
	return e
}

// GetData get this error pointer data.
func (e *Errs) GetData() *Errs {
	return e
}

// NewErrContext create a new `errs.IErrs`.
func NewErrContext() IErrs {
	obj := &Errs{Status: "FAILED"}
	if ErrorDebugLocation == "1" {
		obj.Location = getCallerFunc()
	}
	return obj
}
