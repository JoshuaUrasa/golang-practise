package config

import (
	"errors"
	"os"
)

type Config struct {
	AppEnv           string
	Port             string
	DB               string
	JwtAccessSecret  string
	JwtRefreshSecret string
	BaseURL          string
	UploadDir        string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		Port:             getEnv("PORT", "8080"),
		DB:               os.Getenv("database"),
		JwtAccessSecret:  os.Getenv("access_secret"),
		JwtRefreshSecret: os.Getenv("refresh_secret"),
		BaseURL:          getEnv("BASE_URL", "http://localhost:8000"),
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
	}

	if cfg.DB == "" {
		return nil, errors.New("database is required")
	}

	if cfg.JwtAccessSecret == "" {
		return nil, errors.New("access_secret is required")
	}

	if cfg.JwtRefreshSecret == "" {
		return nil, errors.New("refresh_secret is required")
	}
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
