package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/manuelam2003/triclone/internal/validator"
)

type ExpenseParticipant struct {
	ID         int64     `json:"id"`
	ExpenseID  int64     `json:"expense_id"`
	UserID     int64     `json:"user_id"`
	AmountOwed float64   `json:"amount_owed"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ExpenseParticipantModel struct {
	DB *sql.DB
}

func ValidateParticipant(v *validator.Validator, participant *ExpenseParticipant) {
	v.Check(participant.ExpenseID > 0, "expense_id", "must be non negative")
	v.Check(participant.UserID > 0, "user_id", "must be non negative")
	v.Check(participant.AmountOwed > 0, "amount_owed", "must be non negative")
}

func (m ExpenseParticipantModel) GetAllForGroupAndExpense(groupID, expenseID int64, filters Filters) ([]*ExpenseParticipant, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, expense_id, user_id, amount_owed, updated_at
		FROM expense_participants
		WHERE expense_id = $1 AND expense_id IN (SELECT id FROM expenses WHERE group_id = $2)
		ORDER BY %s %s
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, expenseID, groupID, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var participants []*ExpenseParticipant

	for rows.Next() {
		var participant ExpenseParticipant
		err := rows.Scan(
			&totalRecords,
			&participant.ID,
			&participant.ExpenseID,
			&participant.UserID,
			&participant.AmountOwed,
			&participant.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		participants = append(participants, &participant)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return participants, metadata, nil
}

func (m ExpenseParticipantModel) Insert(participant *ExpenseParticipant) error {
	query := `
		INSERT INTO expense_participants(expense_id, user_id, amount_owed)
		VALUES ($1, $2, $3)
		RETURNING id, updated_at`

	args := []any{participant.ExpenseID, participant.UserID, participant.AmountOwed}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&participant.ID, &participant.UpdatedAt)
}

func (m ExpenseParticipantModel) Get(participantID int64) (*ExpenseParticipant, error) {
	query := `
		SELECT id, expense_id, user_id, amount_owed, updated_at
		FROM expense_participants
		WHERE id = $1`

	var participant ExpenseParticipant

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, participantID).Scan(
		&participant.ID,
		&participant.ExpenseID,
		&participant.UserID,
		&participant.AmountOwed,
		&participant.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &participant, nil
}

func (m ExpenseParticipantModel) Update(participant *ExpenseParticipant) error {
	query := `
	UPDATE expense_participants
	SET amount_owed = $1, updated_at = NOW()
	WHERE id = $2
	RETURNING updated_at`

	args := []any{participant.AmountOwed, participant.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&participant.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m ExpenseParticipantModel) Delete(participantID int64) error {
	query := `
		DELETE FROM expense_participants
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, participantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
