package server

import (
	"expense-tracker/internal/auth"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, accessSecret, refreshSecret string) *echo.Echo {
	e := echo.New()

	//initialize auth service and handler
	jwtService := auth.NewJWTService(accessSecret, refreshSecret)
	authService := auth.NewService(db, jwtService)
	authHandler := auth.NewHandler(authService)

	//Api group
	api := e.Group("/api")

	//version group
	v1 := api.Group("/v1")

	authRoutes := v1.Group("/auth")

	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	return e
}
