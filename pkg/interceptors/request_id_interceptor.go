package interceptors

import (
	"parking-service/pkg/contexts"
	"parking-service/pkg/logs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AttachRequestID ...
func AttachRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {

			id := uuid.New().String()
			ec.Response().Header().Add(echo.HeaderXRequestID, id)

			bearerContext := contexts.EnsureBearerContext(ec)
			bearerContext.RequestID = id
			bearerContext.GetLogger().Upsert(logs.ClientRequestID, id)
			bearerContext.GetLogger().Update()

			return next(bearerContext)
		}
	}
}
