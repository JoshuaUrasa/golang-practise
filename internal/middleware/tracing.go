package middleware

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func Tracing() echo.MiddlewareFunc {
	tracer := otel.Tracer("expense-tracker/http")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()

			spanName := fmt.Sprintf("%s %s", req.Method, req.URL.Path)

			ctx, span := tracer.Start(req.Context(), spanName)
			defer span.End()

			span.SetAttributes(
				attribute.String("http.method", req.Method),
				attribute.String("http.path", req.URL.Path),
			)

			c.SetRequest(req.WithContext(ctx))

			err := next(c)

			if err != nil {
				span.RecordError(err)
			}

			return err
		}
	}
}
