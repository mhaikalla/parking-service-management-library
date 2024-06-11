package errs

import (
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrsSuite struct {
	suite.Suite
}

func (es *ErrsSuite) TestErrsObjectCreation() {
	err := NewErrContext().SetCode(HTTPCLientRequestBodyErr).SetError(errors.New("ERROR_TEST")).SetMessage("ERROR_TEST_MESSAGE")
	es.EqualError(err.OrigError, "ERROR_TEST", "original error must expected")
	es.EqualError(err, "code=212;code_detail=212;status=FAILED;message=ERROR_TEST_MESSAGE;error=ERROR_TEST;location=", "error message must expected")
	es.Error(err.SetLocation())
	es.EqualError(err, "code=212;code_detail=212;status=FAILED;message=ERROR_TEST_MESSAGE;error=ERROR_TEST;location=parking-service/pkg/errs.(*ErrsSuite).TestErrsObjectCreation:19")
	es.NotEmpty(err.GetData(), "error context should not empty")
	es.NotPanics(func() { es.Error(err.Log(log.Println)) }, "on error not panics")
}

func TestUpstreamsRequestSuite(t *testing.T) {
	suite.Run(t, new(ErrsSuite))
}
