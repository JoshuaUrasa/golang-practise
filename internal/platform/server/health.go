package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check godoc
// @Summary      Health check
// @Description  Check API and database health
// @Tags         system
// @Produce      json
// @Success      200  {object}  server.HealthResponse
// @Failure      500  {object}  server.HealthResponse
// @Router       /health [get]
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
