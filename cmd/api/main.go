package main

import (
	"expense-tracker/internal/expense"
	"expense-tracker/internal/platform/config"
	"expense-tracker/internal/platform/database"
	"expense-tracker/internal/platform/logger"
	"expense-tracker/internal/platform/server"
	"expense-tracker/internal/user"
	"fmt"

	"github.com/joho/godotenv"
)

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

	//initialize logger
	appLogger := logger.New(cfg.AppEnv)

	e := server.NewRouter(db, cfg.JwtAccessSecret, cfg.JwtRefreshSecret, appLogger)
	fmt.Println("Starting server on port", cfg.Port)

	err = e.Start(":" + cfg.Port)
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
	}

	appLogger.Info("config loaded", "env", cfg.AppEnv, "port", cfg.Port)
}
