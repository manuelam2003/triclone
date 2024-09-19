package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/manuelam2003/triclone/internal/validator"
)

type Expense struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	PaidBy      *int64    `json:"paid_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpenseModel struct {
	DB *sql.DB
}

func ValidateExpense(v *validator.Validator, expense *Expense) {
	v.Check(expense.GroupID > 0, "group_id", "must be non negative")
	v.Check(expense.Amount > 0.0, "amount", "must be non negative")
	v.Check(expense.Description != "", "description", "must be provided")
	v.Check(len(expense.Description) <= 500, "description", "must not be more than 500 bytes long")
}

func (m ExpenseModel) Insert(expense *Expense) error {
	query := `
		INSERT INTO expenses(group_id, amount, description, paid_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	args := []any{expense.GroupID, expense.Amount, expense.Description, *expense.PaidBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
}
