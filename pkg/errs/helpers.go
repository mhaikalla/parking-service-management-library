package errs

import (
	"regexp"
	"runtime"
	"strconv"
)

var (
	reMaskedMessage = regexp.MustCompile(`^[A-Z0-9_\-]*$`)
)

// BubbleOriErr bubling up deep `*Errs`.
func BubbleOriErr(err *Errs) *Errs {
	if err == nil {
		return err
	}
	// peek the original error
	// return deep `*Errs`
	if e, ok := err.OrigError.(*Errs); ok {
		return BubbleOriErr(e)
	}
	return err
}

// getCallerFunc get caller function.
func getCallerFunc() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + ":" + strconv.Itoa(line)
}

// MaskingError masking an error passing to client.
// If error is not `*Errs` wrapped it with `GeneralErrorHandler`.
// If env `MASKING_ERROR_MESSAGE==1` try to masking it, only when the message not in uppercased-non-space string.
func MaskingError(err error) error {
	if err == nil {
		return err
	}

	isMasking := MaskingErrorMessage == "1"
	var resErr *Errs
	outer, ok := err.(*Errs)

	if ok {
		oriErr := BubbleOriErr(outer)
		resErr = oriErr
	} else {
		resErr = NewErrContext().SetCode(GeneralErrorHandler).SetError(err)
	}

	msg := resErr.Message
	if isMasking && (!reMaskedMessage.MatchString(msg) || msg == "") {
		return resErr.SetMessage(MessageByCode[resErr.Code])
	}

	return resErr
}
