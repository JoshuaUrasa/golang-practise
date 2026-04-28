package expense

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateExpense(ctx context.Context, userId uint, req CreateExpenseRequest) (*Expense, error) {
	tracer := otel.Tracer("expense-service")

	_, span := tracer.Start(ctx, "CreateExpense")
	defer span.End()

	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if req.Category == "" {
		return nil, errors.New("category is required")
	}

	expense := Expense{
		UserID:      userId,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
	}

	if err := s.db.Create(&expense).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	if err := s.db.Preload("User").First(&expense, expense.ID).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &expense, nil
}

func (s *Service) ListAllExpenses(userId uint) ([]Expense, error) {
	var expenses []Expense
	if err := s.db.Where("user_id = ?", userId).Find(&expenses).Error; err != nil {
		return nil, err
	}
	s.db.Preload("User").Where("user_id = ?", userId).Find(&expenses)
	return expenses, nil
}

func (s *Service) GetExpenseById(userId, expenseId uint) (*Expense, error) {
	var expense Expense
	if err := s.db.Where("user_id = ? AND id = ?", userId, expenseId).First(&expense).Error; err != nil {
		return nil, err
	}

	s.db.Preload("User").Where("id = ? AND user_id = ?", expenseId, userId).First(&expense)
	return &expense, nil
}

func (s *Service) Update(userID, expenseID uint, req UpdateExpenseRequest) (*Expense, error) {
	expense, err := s.GetExpenseById(userID, expenseID)
	if err != nil {
		return nil, err
	}

	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if req.Category == "" {
		return nil, errors.New("category is required")
	}

	expense.Amount = req.Amount
	expense.Category = req.Category
	expense.Description = req.Description

	if err := s.db.Save(expense).Error; err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *Service) Delete(userID, expenseID uint) error {
	expense, err := s.GetExpenseById(userID, expenseID)
	if err != nil {
		return err
	}

	return s.db.Delete(expense).Error
}
