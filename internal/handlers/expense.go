package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/Archiker-715/expense-tracker-api/internal/errs"
	"github.com/Archiker-715/expense-tracker-api/internal/repository/pg"
	"github.com/Archiker-715/expense-tracker-api/internal/usecase"
	"github.com/Archiker-715/expense-tracker-api/pkg/httpserver"
)

type ExpenseHandler struct {
	expense *usecase.ExpenseService
}

func NewExpenseHandler(repo *pg.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{expense: usecase.NewExpenseService(repo)}
}

func (e *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	idStr := query.Get("id")
	if idStr == "" {
		expenses, err := e.expense.GetExpenses(ctx)
		if err != nil {
			errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
		}
		if err := httpserver.JsonEncode(w, expenses, 0); err != nil {
			return
		}
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
	}

	expense, err := e.expense.GetExpenseById(ctx, uint(id))
	if err != nil {
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
	}
	if err := httpserver.JsonEncode(w, expense, 0); err != nil {
		return
	}

}

func (e *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense entity.ExpenseCreate
	if err := httpserver.JsonDecode(w, r, &expense, 0); err != nil {
		return
	}

	ctx := r.Context()
	newExpenseId, err := e.expense.CreateExpense(ctx, expense)
	if err != nil {
		if errors.Is(err, usecase.ParseTimeError) {
			errs.WriteError(w, 0, http.StatusBadRequest, fmt.Sprint(err))
		}
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
	}

	if err := httpserver.JsonEncode(w, newExpenseId, 0); err != nil {
		return
	}
}
