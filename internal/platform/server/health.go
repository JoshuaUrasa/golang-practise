package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *echo.Context) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "database error",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "database down",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
