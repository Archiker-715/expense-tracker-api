package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Archiker-715/expense-tracker-api/internal/auth"
	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/Archiker-715/expense-tracker-api/internal/repository/pg"
	"github.com/google/uuid"
)

type ExpenseService struct {
	repo *pg.ExpenseRepository
}

func NewExpenseService(repo *pg.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

var ParseTimeError error = fmt.Errorf("Parsing date error. Enter date as %v", time.DateOnly)

func (e *ExpenseService) GetExpenses(ctx context.Context) ([]entity.Expense, error) {
	userId, ok := auth.UserFromContext(ctx)
	if !ok {
		return []entity.Expense{}, errors.New("empty userId in ctx")
	}
	return e.repo.GetExpenses(userId)
}

func (e *ExpenseService) GetExpenseById(ctx context.Context, id int) (entity.Expense, error) {
	_, exp, err := e.checkExpOwner(ctx, id)
	if err != nil {
		return entity.Expense{}, err
	}
	return exp, nil
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

	return e.repo.CreateExpense(&expense)
}

func (e *ExpenseService) UpdateExpense(ctx context.Context, expense entity.ExpenseUpdate, expenseId int) error {
	if expense.Date != nil {
		if _, err := time.Parse(time.DateOnly, (*expense.Date)); err != nil {
			return ParseTimeError
		}
	}

	userId, _, err := e.checkExpOwner(ctx, expenseId)
	if userId == uuid.Nil || err != nil {
		return fmt.Errorf("usersId check: %w", err)
	}

	var exp entity.Expense
	exp.ID = uint(expenseId)
	if expense.Date != nil {
		exp.Date = *expense.Date
	}
	if expense.Amount != nil {
		exp.Amount = *expense.Amount
	}
	if expense.Category != nil {
		exp.Category = *expense.Category
	}

	return e.repo.UpdateExpense(&exp)
}

func (e *ExpenseService) DeleteExpense(ctx context.Context, expenseId int) error {
	userId, _, err := e.checkExpOwner(ctx, expenseId)
	if userId == uuid.Nil || err != nil {
		return fmt.Errorf("usersId check: %w", err)
	}
	return e.repo.DeleteExpense(uint(expenseId))
}

func (e *ExpenseService) checkExpOwner(ctx context.Context, expenseId int) (userId uuid.UUID, exp entity.Expense, err error) {
	userId, ok := auth.UserFromContext(ctx)
	if !ok {
		return uuid.Nil, entity.Expense{}, errors.New("empty userId in ctx")
	}

	expFromDB, err := e.repo.GetExpenseById(uint(expenseId))
	if err != nil {
		return uuid.Nil, entity.Expense{}, err
	}

	if expFromDB.InsertedBy != userId {
		return uuid.Nil, entity.Expense{}, errors.New("not enough rights for action")
	}

	return userId, expFromDB, nil
}
