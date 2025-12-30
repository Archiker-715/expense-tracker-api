package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Archiker-715/expense-tracker-api/internal/auth"
	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/Archiker-715/expense-tracker-api/internal/repository/pg"
)

type ExpenseService struct {
	repo *pg.ExpenseRepository
}

func NewExpenseService(repo *pg.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

var ParseTimeError error = fmt.Errorf("Parsing date error. Enter date as %v", time.DateOnly)

func (e *ExpenseService) GetExpenses(ctx context.Context) ([]entity.Expense, error) {
	return e.repo.GetExpenses()
}

func (e *ExpenseService) GetExpenseById(ctx context.Context, id uint) (entity.Expense, error) {
	return e.repo.GetExpenseById(id)
}

func (e *ExpenseService) CreateExpense(ctx context.Context, newExpense entity.ExpenseCreate) (entity.ID, error) {

	if _, err := time.Parse(time.DateOnly, newExpense.Date); err != nil {
		return entity.ID{}, ParseTimeError
	}

	userId, ok := auth.UserFromContext(ctx)
	if !ok {
		return entity.ID{}, errors.New("empty userId in ctx")
	}

	expense := entity.Expense{
		Date:       newExpense.Date,
		Amount:     newExpense.Amount,
		Category:   newExpense.Category,
		InsertedBy: userId,
	}

	newExpenseId, err := e.repo.CreateExpense(&expense)
	if err != nil {
		return entity.ID{}, fmt.Errorf("create expense error: %w", err)
	}
	return newExpenseId, nil
}
