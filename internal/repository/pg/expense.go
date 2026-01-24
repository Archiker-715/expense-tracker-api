package pg

import (
	"time"

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
	if err = e.DB.Where("inserted_by = ?", userId).Find(&expenses).Error; err != nil {
		return
	}
	return
}

func (e *ExpenseRepository) GetExpensesByDateInterval(userId uuid.UUID, expenseId uint, startDate, endDate time.Time) (expenses []entity.Expense, err error) {
	if expenseId != 0 {
		if err = e.DB.Where("inserted_by = ?", userId).Where("inserted >= ? AND inserted <= ?", startDate, endDate).Where("ID = ?", expenseId).Find(&expenses).Error; err != nil {
			return
		}
	} else {
		if err = e.DB.Where("inserted_by = ?", userId).Where("inserted >= ? AND inserted <= ?", startDate, endDate).Find(&expenses).Error; err != nil {
			return
		}
	}
	return
}

func (e *ExpenseRepository) GetExpensesByPastDate(userId uuid.UUID, expenseId uint, pastDate time.Time) (expenses []entity.Expense, err error) {
	if expenseId != 0 {
		if err = e.DB.Where("inserted_by = ?", userId).Where("inserted >= ? AND inserted <= ?", pastDate, time.Now()).Where("ID = ?", expenseId).Find(&expenses).Error; err != nil {
			return
		}
	} else {
		if err = e.DB.Where("inserted_by = ?", userId).Where("inserted >= ? AND inserted <= ?", pastDate, time.Now()).Find(&expenses).Error; err != nil {
			return
		}
	}
	return
}

func (e *ExpenseRepository) GetExpenseById(userId uuid.UUID, expenseId uint) (expense entity.Expense, err error) {
	if err = e.DB.Where("inserted_by = ?", userId).Where("ID = ?", expenseId).Find(&expense).Error; err != nil {
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

func (e *ExpenseRepository) UpdateExpense(userId uuid.UUID, expense *entity.Expense) error {
	if err := e.DB.Where("inserted_by = ?", userId).Updates(expense).Error; err != nil {
		return err
	}
	return nil
}

func (e *ExpenseRepository) DeleteExpense(userId uuid.UUID, id uint) error {
	if err := e.DB.Where("inserted_by = ?", userId).Delete(entity.Expense{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
