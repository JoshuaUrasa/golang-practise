package middleware

import (
	"expense-tracker/internal/auth"
	"strings"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService *auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return c.JSON(401, map[string]string{
					"error": "missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(401, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			claims, err := jwtService.ValidateAccessToken(parts[1])

			if err != nil {
				return c.JSON(401, map[string]string{
					"error": "invalid or expired token",
				})

			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
