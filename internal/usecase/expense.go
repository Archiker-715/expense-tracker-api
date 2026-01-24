package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
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

var availableCategories entity.AvailableExpCategories

var ParseTimeError error = fmt.Errorf("Parsing date error. Enter date as %v", time.DateOnly)

func (e *ExpenseService) GetExpenses(ctx context.Context, expenseId int, dateFilter entity.DateFilter) ([]entity.Expense, error) {
	userId, ok := auth.UserFromContext(ctx)
	if !ok {
		return []entity.Expense{}, errors.New("empty userId in ctx")
	}

	if dateFilter.PastDate != "" {
		if dateFilter.PastDate != "week" && dateFilter.PastDate != "month" && dateFilter.PastDate != "3 months" {
			return []entity.Expense{}, errors.New("validation date failed: must compile: week or month or 3 months")
		}

		if !dateFilter.StartDate.IsZero() || !dateFilter.EndDate.IsZero() {
			return []entity.Expense{}, errors.New("validation date filters failed: must have only pastDate or only startDate and endDate")
		}

		var date time.Time
		date = time.Now()
		switch dateFilter.PastDate {
		case "week":
			date = date.AddDate(0, 0, -7)
		case "month":
			date = date.AddDate(0, 0, -30)
		case "3 months":
			date = date.AddDate(0, 0, -90)
		}
		return e.repo.GetExpensesByPastDate(userId, uint(expenseId), date)
	} else {
		if !dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
			return e.repo.GetExpensesByDateInterval(userId, uint(expenseId), dateFilter.StartDate, dateFilter.EndDate)
		} else if dateFilter.StartDate.IsZero() && !dateFilter.EndDate.IsZero() {
			return []entity.Expense{}, errors.New("validation failed: startDate cannot be empty")
		} else if !dateFilter.StartDate.IsZero() && dateFilter.EndDate.IsZero() {
			return []entity.Expense{}, errors.New("validation failed: endDate cannot be empty")
		}
	}
	if expenseId <= 0 {
		return e.repo.GetExpenses(userId)
	} else {
		exps := make([]entity.Expense, 0, 1)
		exp, err := e.repo.GetExpenseById(userId, uint(expenseId))
		if err != nil {
			return []entity.Expense{}, err
		}
		if exp.ID == 0 {
			return []entity.Expense{}, errors.New("expense not found")
		}
		exps = append(exps, exp)
		return exps, nil
	}
}

func (e *ExpenseService) CreateExpense(ctx context.Context, newExpense entity.ExpenseCreate) (entity.ID, error) {
	if _, err := time.Parse(time.DateOnly, newExpense.Date); err != nil {
		return entity.ID{}, ParseTimeError
	}

	nexExpCat, ok := categoryAvailable(newExpense.Category)
	if !ok {
		return entity.ID{}, notAvailableCategoryError()
	}

	userId, ok := auth.UserFromContext(ctx)
	if !ok {
		return entity.ID{}, errors.New("empty userId in ctx")
	}

	expense := entity.Expense{
		Date:       newExpense.Date,
		Amount:     newExpense.Amount,
		Category:   nexExpCat,
		InsertedBy: userId,
		Inserted:   time.Now(),
	}

	return e.repo.CreateExpense(&expense)
}

func (e *ExpenseService) UpdateExpense(ctx context.Context, expense entity.ExpenseUpdate, expenseId int) error {
	if expense.Date != nil {
		if _, err := time.Parse(time.DateOnly, (*expense.Date)); err != nil {
			return ParseTimeError
		}
	}

	userId, ok := auth.UserFromContext(ctx)
	if !ok {
		return errors.New("empty userId in ctx")
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
		nexExpCat, ok := categoryAvailable((*expense.Category))
		if !ok {
			return notAvailableCategoryError()
		}
		exp.Category = nexExpCat
	}

	return e.repo.UpdateExpense(userId, &exp)
}

func (e *ExpenseService) DeleteExpense(ctx context.Context, expenseId int) error {
	userId, ok := auth.UserFromContext(ctx)
	if !ok {
		return errors.New("empty userId in ctx")
	}
	return e.repo.DeleteExpense(userId, uint(expenseId))
}

func categoryAvailable(category string) (newExpCat string, available bool) {
	availableCategories.SetAvailableExpCategories()
	newExpCat = strings.ToLower(category)
	available = true
	if !slices.Contains(availableCategories.AvailableCategories, newExpCat) {
		return "", false
	}
	newExpCat = strings.ToUpper(newExpCat[:1]) + strings.ToLower(newExpCat[1:])
	return
}

func notAvailableCategoryError() error {
	return fmt.Errorf("Validation category failed. Available is %v", availableCategories.AvailableCategories)
}
