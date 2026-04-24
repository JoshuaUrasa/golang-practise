package expense

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateExpense(userId uint, req CreateExpenseRequest) (*Expense, error) {
	expense := Expense{
		UserID:      userId,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
	}

	if err := s.db.Create(&expense).Error; err != nil {
		return nil, err

	}
	if err := s.db.Preload("User").First(&expense, expense.ID).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (s *Service) ListAllExpenses(userId uint) ([]Expense, error) {
	var expenses []Expense
	if err := s.db.Where("UserID = ?", userId).Find(&expenses).Error; err != nil {
		return nil, err
	}
	s.db.Preload("User").Where("user_id = ?", userId).Find(&expenses)
	return expenses, nil
}

func (s *Service) GetExpenseById(userId, expenseId uint) (*Expense, error) {
	var expense Expense
	if err := s.db.Where("UserID = ? AND ID = ?", userId, expenseId).First(&expense).Error; err != nil {
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
