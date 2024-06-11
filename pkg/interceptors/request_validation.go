package interceptors

import (
	"parking-service/pkg/contexts"

	"github.com/labstack/echo/v4"
)

// SetRequestValidationData ...
func SetRequestValidationData() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			bearerContext := contexts.EnsureBearerContext(ec)
			bearerContext.RequestValidationData.UserAgent = ec.Request().Header.Get("user-agent")
			bearerContext.RequestValidationData.XHV = ec.Request().Header.Get("x-hv")
			bearerContext.RequestValidationData.XSignature = ec.Request().Header.Get("x-signature")
			bearerContext.RequestValidationData.XVersionApp = ec.Request().Header.Get("x-version-app")
			bearerContext.RequestValidationData.RequestMethod = ec.Request().Method
			bearerContext.RequestValidationData.RequestPath = ec.Request().URL.Path

			return next(bearerContext)
		}
	}
}
