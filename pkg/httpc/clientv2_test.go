package httpc

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UpstreamsRequestSuite struct {
	suite.Suite
}

// respFunc send response using supllied args
func respFunc(rw http.ResponseWriter, code int, body []byte, contentType string) error {
	switch contentType {
	case "json":
		rw.Header().Add("content-type", "application/json")
	case "text":
		rw.Header().Add("content-type", "plain/text")
	default:
		rw.Header().Add("content-type", "plain/html")
	}
	rw.WriteHeader(code)
	if _, err := rw.Write(body); err != nil {
		return err
	}
	return nil
}

func (urs *UpstreamsRequestSuite) TestSuccessOnGetMethod() {

	handler := func(rw http.ResponseWriter, req *http.Request) {
		testData := `{}`
		err := respFunc(rw, 200, []byte(testData), "json")
		urs.NoError(err, "Handling request should no error")
	}

	server := httptest.NewServer(http.HandlerFunc(handler))

	defer server.Close()

	client := &http.Client{Timeout: time.Second * 5}

	req := UpstreamsRequest{
		URL:          server.URL,
		Method:       "GET",
		Client:       client,
		SuccessCodes: []int{200},
		Headers: map[string]string{
			"X-API-KEY": "FOOOBARR123",
		},
		Bearer: "FOOBAZ12345",
	}

	res := req.Request()

	urs.Nil(res.GetError(), "should no error")
	urs.Equal(200, res.GetCode(), "response code must expected")
	urs.NotEmpty(res.GetHeaders(), "headers not empty")
	urs.True(res.IsSuccess(), "response must success")
}

func (urs *UpstreamsRequestSuite) TestErrorOnParsingPayload() {

	handler := func(rw http.ResponseWriter, req *http.Request) {
		testData := `{}`
		err := respFunc(rw, 200, []byte(testData), "json")
		urs.NoError(err, "Handling request should no error")
	}

	server := httptest.NewServer(http.HandlerFunc(handler))

	defer server.Close()

	client := &http.Client{Timeout: time.Second * 5}

	req := UpstreamsRequest{
		URL:          server.URL,
		Method:       "POST",
		Client:       client,
		SuccessCodes: []int{200},
		BodyPayload:  math.Inf(-1),
	}

	res := req.Request()

	urs.NotNil(res.GetError(), "should error on parsing body payload")
	urs.False(res.IsSuccess(), "should no success")
}

func TestUpstreamsRequestSuite(t *testing.T) {
	suite.Run(t, new(UpstreamsRequestSuite))
}

func Test_maskingBearer(t *testing.T) {
	type args struct {
		bearer string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"odd",
			args{bearer: "Bearer jjsks0ksjss88484jsks"},
			"Bearer jjsks0ks****8484jsks",
		},
		{
			"even",
			args{bearer: "Bearer jjsks0ksjss88484jsk"},
			"Bearer jjsks0ks***88484jsk",
		},
		{
			"5 digit",
			args{bearer: "Bearer jjsks"},
			"Bearer jjsks",
		},
		{
			"6 digit",
			args{bearer: "Bearer jjsks0"},
			"Bearer j****0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, maskingBearer(tt.args.bearer), tt.name)
		})
	}
}
