package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/manuelam2003/triclone/internal/validator"
)

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidateGroup(v *validator.Validator, group *Group) {
	v.Check(group.Name != "", "name", "must be provided")
	v.Check(len(group.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(group.CreatedBy != 0, "created_by", "must be provided")
	v.Check(group.CreatedBy > 0, "created_by", "must be a positive integer")
}

type GroupModel struct {
	DB *sql.DB
}

func (g GroupModel) Insert(group *Group) error {
	query := `
		INSERT INTO groups (name, created_by)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	args := []any{group.Name, group.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return g.DB.QueryRowContext(ctx, query, args...).Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
}

func (g GroupModel) Get(id int64) (*Group, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, created_by, created_at, updated_at
		FROM groups
		WHERE id = $1`

	var group Group

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := g.DB.QueryRowContext(ctx, query, id).Scan(
		&group.ID,
		&group.Name,
		&group.CreatedBy,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &group, nil
}

func (g GroupModel) Update(group *Group) error {
	query := `
		UPDATE groups
		SET name = $1, created_by = $2, updated_at = NOW()
		WHERE id = $3 AND updated_at = $4
		RETURNING updated_at`

	args := []any{group.Name, group.CreatedBy, group.ID, group.UpdatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := g.DB.QueryRowContext(ctx, query, args...).Scan(&group.UpdatedAt)
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

func (g GroupModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM groups
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := g.DB.ExecContext(ctx, query, id)
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

func (g GroupModel) GetAll(name string, createdBy int64, filters Filters) ([]*Group, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, created_by, created_at, updated_at
		FROM groups
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (created_by = $2 OR $2 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, createdBy, filters.limit(), filters.offset()}

	rows, err := g.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	groups := []*Group{}

	for rows.Next() {
		var group Group

		err := rows.Scan(
			&totalRecords,
			&group.ID,
			&group.Name,
			&group.CreatedBy,
			&group.CreatedAt,
			&group.UpdatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return groups, metadata, nil
}
