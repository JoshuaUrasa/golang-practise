package middleware

import (
	"errors"
	"expense-tracker/internal/platform/metrics"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

func Metrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Implementation for metrics middleware
			start := time.Now()

			err := next(c)

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

			duration := time.Since(start).Seconds()
			statusStr := strconv.Itoa(status)
			method := c.Request().Method
			path := c.Request().URL.Path

			metrics.HTTPRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(method, path, statusStr).Observe(duration)

			return err
		}
	}
}
