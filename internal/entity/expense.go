package entity

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID         uint      `json:"expenseId" gorm:"primaryKey;autoIncrement"`
	Date       string    `json:"expenseDate"`
	Amount     int       `json:"expenseAmount" validate:"required" gorm:"check:price >=0"`
	Category   string    `json:"expenseCategory" validate:"required,oneof=Groceries Leisure Electronics Utilities Clothing Health Others"`
	InsertedBy uuid.UUID `json:"inserted_by"`
	Inserted   time.Time `json:"inserted"`
}

type ExpenseCreate struct {
	// ID       uint   `swaggerignore:"true"`
	Date     string `json:"expenseDate,omitempty"`
	Amount   int    `json:"expenseAmount"`
	Category string `json:"expenseCategory" validate:"required,oneof=Groceries Leisure Electronics Utilities Clothing Health Others"`
}

type ExpenseUpdate struct {
	Date     *string `json:"expenseDate,omitempty"`
	Amount   *int    `json:"expenseAmount,omitempty"`
	Category *string `json:"expenseCategory,omitempty" validate:"required,oneof=Groceries Leisure Electronics Utilities Clothing Health Others"`
}

type DateFilter struct {
	PastDate  string
	StartDate time.Time
	EndDate   time.Time
}
