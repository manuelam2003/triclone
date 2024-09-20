package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/manuelam2003/triclone/internal/validator"
)

type Settlement struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	PayerID   *int64    `json:"payer_id"`
	PayeeID   *int64    `json:"payee_id"`
	Amount    float64   `json:"amount"`
	SettledAt time.Time `json:"settled_at"`
}

func ValidateSettlement(v *validator.Validator, settlement *Settlement) {
	v.Check(settlement.GroupID > 0, "group_id", "must be positive")

	if settlement.PayerID != nil {
		v.Check(*settlement.PayerID > 0, "payer_id", "must be positive")
	}

	if settlement.PayeeID != nil {
		v.Check(*settlement.PayeeID > 0, "payee_id", "must be positive")
	}

	if settlement.PayeeID != nil && settlement.PayerID != nil {
		v.Check(*settlement.PayeeID != *settlement.PayerID, "payee_id", "payer and payee cannot be the same")
	}

	v.Check(settlement.Amount > 0, "amount", "must be positive")
}

type SettlementModel struct {
	DB *sql.DB
}

func (m SettlementModel) GetAllForGroup(groupID int64, filters Filters) ([]*Settlement, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, group_id, payer_id, payee_id, amount, settled_at
		FROM settlements
		WHERE group_id = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, groupID, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	var settlements []*Settlement

	for rows.Next() {
		var settlement Settlement
		err := rows.Scan(
			&totalRecords,
			&settlement.ID,
			&settlement.GroupID,
			&settlement.PayerID,
			&settlement.PayeeID,
			&settlement.Amount,
			&settlement.SettledAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		settlements = append(settlements, &settlement)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return settlements, metadata, nil
}

func (m SettlementModel) Get(settlementID int64, groupID int64) (*Settlement, error) {
	query := `
		SELECT id, group_id, payer_id, payee_id, amount, settled_at
		FROM settlements
		WHERE id = $1 AND group_id = $2`

	var settlement Settlement

	err := m.DB.QueryRow(query, settlementID, groupID).Scan(
		&settlement.ID,
		&settlement.GroupID,
		&settlement.PayerID,
		&settlement.PayeeID,
		&settlement.Amount,
		&settlement.SettledAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &settlement, nil
}

func (m SettlementModel) Insert(settlement *Settlement) error {
	query := `
		INSERT INTO settlements (group_id, payer_id, payee_id, amount, settled_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, settled_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, and scan the returned `id` and `settled_at` fields into the settlement object
	err := m.DB.QueryRowContext(ctx, query, settlement.GroupID, settlement.PayerID, settlement.PayeeID, settlement.Amount, settlement.SettledAt).Scan(
		&settlement.ID, &settlement.SettledAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m SettlementModel) Delete(groupID, settlementID int64) error {
	query := `
		DELETE FROM settlements
		WHERE id = $1 AND group_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, settlementID, groupID)
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
