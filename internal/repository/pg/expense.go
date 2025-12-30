package pg

import (
	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	DB *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{DB: db}
}

func (e *ExpenseRepository) GetExpenses() (expenses []entity.Expense, err error) {
	if err = e.DB.Find(&expenses).Error; err != nil {
		return
	}
	return
}

func (e *ExpenseRepository) GetExpenseById(id uint) (expense entity.Expense, err error) {
	if err = e.DB.Find(&expense, id).Error; err != nil {
		return
	}
	return
}

func (e *ExpenseRepository) CreateExpense(expense *entity.Expense) (entity.ID, error) {
	if err := e.DB.Create(expense).Error; err != nil {
		return entity.ID{}, err
	}
	return entity.ID{ID: int(expense.ID)}, nil
}
