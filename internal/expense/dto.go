package expense

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
