package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (m ExpenseModel) Get(groupID, expenseID int64) (*Expense, error) {
	if expenseID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, group_id, amount, description, paid_by, created_at, updated_at
		FROM expenses
		WHERE id = $1 AND group_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var expense Expense

	err := m.DB.QueryRowContext(ctx, query, expenseID, groupID).Scan(
		&expense.ID,
		&expense.GroupID,
		&expense.Amount,
		&expense.Description,
		&expense.PaidBy,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &expense, nil
}

func (m ExpenseModel) GetAll(groupID int64, description string, paidBy int64, filters Filters) ([]*Expense, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, group_id, amount, description, paid_by, created_at, updated_at
	FROM expenses
	WHERE group_id = $1
	AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
	AND (paid_by = $3 OR $3 = 0)
	ORDER BY %s %s, id ASC
	LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{groupID, description, paidBy, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	expenses := []*Expense{}

	for rows.Next() {
		var expense Expense

		err := rows.Scan(
			&totalRecords,
			&expense.ID,
			&expense.GroupID,
			&expense.Amount,
			&expense.Description,
			&expense.PaidBy,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		expenses = append(expenses, &expense)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return expenses, metadata, nil
}
