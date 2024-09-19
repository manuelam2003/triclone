package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
