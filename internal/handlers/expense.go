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
	"github.com/gorilla/mux"
)

type ExpenseHandler struct {
	expense *usecase.ExpenseService
}

func NewExpenseHandler(repo *pg.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{expense: usecase.NewExpenseService(repo)}
}

// TODO: add date filters: past week, month, last 3 months, custom
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
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
		}

		expense, err := e.expense.GetExpenseById(ctx, id)
		if err != nil {
			errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
		}
		if err := httpserver.JsonEncode(w, expense, 0); err != nil {
			return
		}
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

func (e *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed convert id: %v", err), http.StatusBadRequest)
	}

	var expense entity.ExpenseUpdate
	if err := httpserver.JsonDecode(w, r, &expense, 0); err != nil {
		return
	}

	ctx := r.Context()
	err = e.expense.UpdateExpense(ctx, expense, id)
	if err != nil {
		if errors.Is(err, usecase.ParseTimeError) {
			errs.WriteError(w, 0, http.StatusBadRequest, fmt.Sprint(err))
		}
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
	}
}

func (e *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed convert id: %v", err), http.StatusBadRequest)
	}

	ctx := r.Context()
	err = e.expense.DeleteExpense(ctx, id)
	if err != nil {
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("http error: %v", err))
	}
}
