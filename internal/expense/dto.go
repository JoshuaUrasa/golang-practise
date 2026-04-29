package expense

import "time"

type CreateExpenseRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
}

type UpdateExpenseRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
}

type ExpenseResponse struct {
	ID          uint      `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type DeleteExpenseResponse struct {
	Message string `json:"message"`
}
