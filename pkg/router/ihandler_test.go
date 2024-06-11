package router

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerIfaceTestSuite struct {
	suite.Suite
}

type routerContext struct {
	method  string
	path    string
	handler func(interface{}) error
}

func (c *routerContext) Handle(method, path string, handler func(interface{}) error) {
	c.method = method
	c.path = path
	c.handler = handler
}

func (s *HandlerIfaceTestSuite) TestImplementIface() {
	ctx := routerContext{}

	hNoError := func(interface{}) error {
		return nil
	}

	hError := func(interface{}) error {
		return errors.New("error test")
	}

	ctx.Handle("get", "/test", hNoError)
	s.Equal(ctx.method, "get", "method must set")
	s.Equal(ctx.path, "/test", "path must set")
	s.NoError(ctx.handler(nil), "handler not return error")

	ctx.Handle("post", "/test", hError)
	s.Equal(ctx.method, "post", "method must set")
	s.Equal(ctx.path, "/test", "path must set")
	s.Error(ctx.handler(nil), "handler must return error")

}

func TestHandlerIfaceSuite(t *testing.T) {
	suite.Run(t, new(HandlerIfaceTestSuite))
}
