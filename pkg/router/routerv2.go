package router

import (
	"github.com/mhaikalla/parking-service-management-library/pkg/interceptors"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EchoServerV2 version 2 of EchoServer context
type EchoServerV2 struct {
	server          *echo.Echo
	config          map[string]map[string]interface{}
	nonAuthHandlers [][]interface{}
	authHandlers    [][]interface{}
}

// ServerV2 version 2 of Server interface
type ServerV2 interface {
	IHandleRegisterV2
	GetServer() *echo.Echo
}

// NewEchoServerV2 version 2 of NewEchoServer
func NewEchoServerV2(config map[string]map[string]interface{}) ServerV2 {
	server := echo.New()
	server.HTTPErrorHandler = JSONErrorHandler
	return &EchoServerV2{
		server:          server,
		config:          config,
		nonAuthHandlers: [][]interface{}{},
		authHandlers:    [][]interface{}{},
	}
}

// Handle method to add route config
func (ctx *EchoServerV2) Handle(method, path string, handler func(interface{}) error) {
	handerfunc := func(c echo.Context) error { return handler(c) }
	ctx.nonAuthHandlers = append(ctx.nonAuthHandlers, []interface{}{method, path, handerfunc})
}

// HandleAuth method to add route config but must be authenticated
func (ctx *EchoServerV2) HandleAuth(method, path string, handler func(interface{}) error) {
	handerfunc := func(c echo.Context) error { return handler(c) }
	ctx.authHandlers = append(ctx.authHandlers, []interface{}{method, path, handerfunc})
}

// GetServer function returning echo server
func (ctx *EchoServerV2) GetServer() *echo.Echo {
	server := ctx.server
	conf := ctx.config

	if serverConf, ok := conf["server"]; ok {
		d := serverConf["debug"].(bool)
		server.Debug = d
		server.HideBanner = true

		if middlewaresConf, ok := serverConf["middlewares"]; ok {
			m := middlewaresConf.([]interface{})
			for _, am := range m {
				switch am.(string) {

				case "attach_request_id":
					server.Use(interceptors.AttachRequestID())

				case "log":
					server.Use(middleware.Logger())

				case "jwt":
					initJWTMiddleware(server, conf, ctx.nonAuthHandlers)

				case "requestid":
					server.Use(middleware.RequestID())

				case "recover":
					server.Use(middleware.Recover())
				}

			}
		}
	}

	server.Use(interceptors.SetRequestValidationData())
	iterateRoutes(server, ctx.nonAuthHandlers)
	iterateRoutes(server, ctx.authHandlers)

	return ctx.server
}
