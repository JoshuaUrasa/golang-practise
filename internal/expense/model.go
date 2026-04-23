package expense

import (
	"expense-tracker/internal/user"
	"time"
)

type Expense struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Category    string    `gorm:"not null" json:"category"`
	Description string    `json:"description"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
