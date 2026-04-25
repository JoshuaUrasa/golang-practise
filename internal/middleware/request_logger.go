package middleware

import (
	"errors"
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

func RequestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start)

			// 1. Kwa v5, unwrap response kwanza ili upate struct yenye .Status
			res, unwrapErr := echo.UnwrapResponse(c.Response())
			status := 0
			if unwrapErr == nil {
				status = res.Status
			}

			// 2. Kama kuna error kutoka kwa handler, vuta status code humo
			if err != nil {
				var sc echo.HTTPStatusCoder
				if errors.As(err, &sc) {
					status = sc.StatusCode()
				} else {
					status = 500
				}
			}

			logger.Info("http request",
				"request_id", c.Get("request_id"),
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", status,
				"duration_ms", duration.Milliseconds(),
			)

			return err
		}
	}
}
