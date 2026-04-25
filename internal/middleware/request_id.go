package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			requestId := uuid.New().String()
			c.Set("request_id", requestId)
			c.Response().Header().Set("X-Request-ID", requestId)
			return next(c)
		}
	}
}
