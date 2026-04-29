package expense

import (
	"net/http"
	"strconv"

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

// CreateExpense godoc
// @Summary      Create expense
// @Description  Create a new expense for the authenticated user
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      expense.CreateExpenseRequest  true  "Expense data"
// @Success      201      {object}  expense.ExpenseResponse
// @Failure      400      {object}  expense.ErrorResponse
// @Failure      401      {object}  expense.ErrorResponse
// @Failure      500      {object}  expense.ErrorResponse
// @Router       /api/v1/expenses [post]
func (h *Handler) CreateExpense(c *echo.Context) error {
	// Implementation for creating an expense
	UserId, ok := c.Get("user_id").(uint)

	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var req CreateExpenseRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	expense, err := h.service.CreateExpense(c.Request().Context(), UserId, req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, expense)
}

// ListExpenses godoc
// @Summary      List expenses
// @Description  List all expenses for the authenticated user
// @Tags         expenses
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   expense.ExpenseResponse
// @Failure      401  {object}  expense.ErrorResponse
// @Failure      500  {object}  expense.ErrorResponse
// @Router       /api/v1/expenses [get]
func (h *Handler) ListExpenses(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	expenses, err := h.service.ListAllExpenses(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, expenses)
}

// GetExpenseByID godoc
// @Summary      Get expense
// @Description  Get one expense by ID for the authenticated user
// @Tags         expenses
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Expense ID"
// @Success      200  {object}  expense.ExpenseResponse
// @Failure      400  {object}  expense.ErrorResponse
// @Failure      401  {object}  expense.ErrorResponse
// @Failure      404  {object}  expense.ErrorResponse
// @Router       /api/v1/expenses/{id} [get]
func (h *Handler) GetExpenseByID(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	expenseID, err := parseID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid expense id",
		})
	}

	expense, err := h.service.GetExpenseById(userID, expenseID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, expense)
}

// UpdateExpense godoc
// @Summary      Update expense
// @Description  Update one expense by ID for the authenticated user
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                           true  "Expense ID"
// @Param        request  body      expense.UpdateExpenseRequest  true  "Expense data"
// @Success      200      {object}  expense.ExpenseResponse
// @Failure      400      {object}  expense.ErrorResponse
// @Failure      401      {object}  expense.ErrorResponse
// @Failure      404      {object}  expense.ErrorResponse
// @Router       /api/v1/expenses/{id} [put]
func (h *Handler) UpdateExpense(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	expenseID, err := parseID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid expense id",
		})
	}

	var req UpdateExpenseRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	expense, err := h.service.Update(userID, expenseID, req)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, expense)
}

// DeleteExpense godoc
// @Summary      Delete expense
// @Description  Delete one expense by ID for the authenticated user
// @Tags         expenses
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Expense ID"
// @Success      200  {object}  expense.DeleteExpenseResponse
// @Failure      400  {object}  expense.ErrorResponse
// @Failure      401  {object}  expense.ErrorResponse
// @Failure      404  {object}  expense.ErrorResponse
// @Router       /api/v1/expenses/{id} [delete]
func (h *Handler) DeleteExpense(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	expenseID, err := parseID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid expense id",
		})
	}

	if err := h.service.Delete(userID, expenseID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "expense deleted successfully",
	})
}

func parseID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
