package auth

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register godoc
// @Summary      User Registration
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.RegisterRequest  true  "User Registration Data"
// @Success      201      {object}  auth.AuthResponse
// @Failure      400      {object}  auth.ErrorResponse
// @Failure      500      {object}  auth.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(c *echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	res, err := h.service.Register(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res)
}

// Login godoc
// @Summary      User Login
// @Description  Authenticate a user and return tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.LoginRequest  true  "Login Credentials"
// @Success      200      {object}  auth.AuthResponse
// @Failure      400      {object}  auth.ErrorResponse
// @Failure      401      {object}  auth.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *Handler) Login(c *echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	res, err := h.service.Login(req)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Issue a new access token using a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.RefreshTokenRequest  true  "Refresh token"
// @Success      200      {object}  auth.AuthResponse
// @Failure      400      {object}  auth.ErrorResponse
// @Failure      401      {object}  auth.ErrorResponse
// @Router       /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *echo.Context) error {
	var req RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	res, err := h.service.RefreshToken(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}
