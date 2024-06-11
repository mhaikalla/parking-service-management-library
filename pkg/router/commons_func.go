package router

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
	"github.com/mhaikalla/parking-service-management-library/pkg/logs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Binding for public API.
var (
	CreateSkippedHandler = createSkippedHandler
	InitJWTMiddleware    = initJWTMiddleware

	reSkippedRouteCheck = regexp.MustCompile(`^.*/api/v[1-9]/.*`)
)

// JSONErrorHandler error handler that return json as response
func JSONErrorHandler(err error, ctx echo.Context) {
	if err != nil {
		bc := contexts.EnsureBearerContext(ctx)
		logger := bc.GetLogger()
		wrapped := errs.MaskingError(err)
		if wrapped.(*errs.Errs).HttpCode != "" {
			bc.Response().Header().Set("status-code", wrapped.(*errs.Errs).HttpCode)
		}
		logs.WhenError("JSON ERROR HANDLER", wrapped, logger)
		logs.WhenError("RESPONSE ERROR", bc.JSON(200, wrapped), logger)
	}
}

// createSkippedHandler create skipper function to bypass some route from middleware
func createSkippedHandler(routes [][]interface{}) func(echo.Context) bool {
	mapRoutes := map[string]map[string]int{}
	for _, r := range routes {
		k := strings.ToUpper(r[0].(string))
		v := r[1].(string)
		if e, f := mapRoutes[k]; f {
			e[v] = 1
		} else {
			mapRoutes[k] = map[string]int{v: 1}
		}
	}

	return func(c echo.Context) bool {

		m := c.Request().Method
		p := c.Request().URL.Path

		if !reSkippedRouteCheck.MatchString(p) {
			return true
		}

		_, isOk := mapRoutes[m][p]

		return isOk
	}
}

// initJWTMiddleware init echo jwt middleware using config
func initJWTMiddleware(server *echo.Echo, config map[string]map[string]interface{}, skippedHandler [][]interface{}) {

	var signingMethod string
	var signingKey string
	var enable bool

	if conf, found := config["jwt"]; found {
		enable = conf["enable"].(bool)
		signingMethod = conf["signing_method"].(string)
	} else {
		panic(errors.New("No configuration key jwt found"))
	}

	if conf, found := config["secrets"]; found {
		if secConf, found := conf["jwt"].(map[string]interface{}); found {
			signingKey = secConf["key"].(string)
		} else {
			panic(errors.New("No configuration key jwt found in secrets entry"))
		}
	} else {
		panic(errors.New("No configuration key secrets found"))
	}

	if enable {
		skipper := createSkippedHandler(skippedHandler)
		jc := middleware.JWTConfig{
			SigningKey:    []byte(signingKey),
			SigningMethod: signingMethod,
			Skipper:       skipper,
		}
		server.Use(middleware.JWTWithConfig(jc))
	}
}

// registerRoute register route to echo server
func registerRoute(server *echo.Echo, method, path string, handler echo.HandlerFunc) {
	m := strings.ToUpper(method)
	switch m {
	case "GET":
		server.GET(path, handler)
	case "POST":
		server.POST(path, handler)
	case "PATCH":
		server.PATCH(path, handler)
	case "DELETE":
		server.DELETE(path, handler)
	case "PUT":
		server.PUT(path, handler)
	default:
		break
	}
}

// iterateRoutes iterater list of routes and register each of route to echo server
func iterateRoutes(server *echo.Echo, routes [][]interface{}) {
	for _, h := range routes {
		m := h[0].(string)
		p := h[1].(string)
		h := h[2].(func(echo.Context) error)
		registerRoute(server, m, p, h)
	}
}
