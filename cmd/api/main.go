package main

import (
	"context"
	"expense-tracker/internal/expense"
	"expense-tracker/internal/platform/config"
	"expense-tracker/internal/platform/database"
	"expense-tracker/internal/platform/logger"
	"expense-tracker/internal/platform/server"
	"expense-tracker/internal/platform/telemetry"
	"expense-tracker/internal/user"
	"fmt"
	"log"
	"net/http"

	_ "expense-tracker/api/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// @title           Expense Tracker API
// @version         1.0
// @description     This is a sample expense tracker server.
// @host            localhost:8000
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	//load env file

	err := godotenv.Load("configs/dev.env")

	if err != nil {
		fmt.Println("no env file found")
	}

	//load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		return
	}

	//connect to database
	db, err := database.NewMySQL(cfg.DB)
	if err != nil {
		fmt.Printf("failed to connect to database: %v\n", err)
		return
	}

	//auto migrate database
	err = db.AutoMigrate(&user.User{}, &expense.Expense{})
	if err != nil {
		fmt.Printf("failed to migrate database: %v\n", err)
		return
	}

	ctx := context.Background()

	tracerProvider, err := telemetry.InitTracer(ctx, "expense-tracker")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}

	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	//initialize logger
	appLogger := logger.New(cfg.AppEnv)
	appLogger.Info("config loaded", "env", cfg.AppEnv, "port", cfg.Port)

	e := server.NewRouter(db, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, appLogger)
	fmt.Println("Starting server on port", cfg.Port)

	//initialize swagger
	e.GET("/swagger", func(c *echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/swagger/", func(c *echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/swagger//index.html", func(c *echo.Context) error {
		target := "/swagger/index.html"
		if c.QueryString() != "" {
			target += "?" + c.QueryString()
		}
		return c.Redirect(http.StatusMovedPermanently, target)
	})
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.URL("/swagger/doc.json")))

	err = e.Start(":" + cfg.Port)
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
	}

}
