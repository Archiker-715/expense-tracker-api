package pg

import (
	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	DB *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{DB: db}
}

func (e *ExpenseRepository) GetExpenses(userId uuid.UUID) (expenses []entity.Expense, err error) {
	if err = e.DB.Where("InsertedBy = ?", userId).Find(&expenses).Error; err != nil {
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

func (e *ExpenseRepository) UpdateExpense(expense *entity.Expense) error {
	if err := e.DB.Save(expense).Error; err != nil {
		return err
	}
	return nil
}

func (e *ExpenseRepository) DeleteExpense(id uint) error {
	if err := e.DB.Delete(entity.Expense{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
