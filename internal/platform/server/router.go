package server

import (
	"expense-tracker/internal/auth"
	"expense-tracker/internal/expense"
	"expense-tracker/internal/middleware"
	"expense-tracker/internal/platform/metrics"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, accessSecret, refreshSecret string, logger *slog.Logger) *echo.Echo {
	e := echo.New()

	//initialize middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.Tracing())
	e.Use(middleware.RequestLogger(logger))
	e.Use(middleware.Metrics())

	//initialize auth service and handler
	jwtService := auth.NewJWTService(accessSecret, refreshSecret)
	authService := auth.NewService(db, jwtService)
	authHandler := auth.NewHandler(authService)

	expenseService := expense.NewService(db)
	expenseHandler := expense.NewHandler(expenseService)

	//metrics initialize
	metrics.Register()
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	//health check endpoint
	healthHandler := NewHealthHandler(db)
	e.GET("/health", healthHandler.Check)

	//Api group
	api := e.Group("/api")

	//version group.
	v1 := api.Group("/v1")

	authRoutes := v1.Group("/auth")

	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/refresh", authHandler.RefreshToken)

	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtService))

	expenseGroup := protected.Group("/expenses")
	expenseGroup.GET("", expenseHandler.ListExpenses)
	expenseGroup.POST("", expenseHandler.CreateExpense)
	expenseGroup.GET("/:id", expenseHandler.GetExpenseByID)
	expenseGroup.PUT("/:id", expenseHandler.UpdateExpense)
	expenseGroup.DELETE("/:id", expenseHandler.DeleteExpense)

	protected.GET("/me", func(c *echo.Context) error {
		return c.JSON(200, map[string]any{
			"user_id": c.Get("user_id"),
			"email":   c.Get("email"),
		})
	})

	return e
}
