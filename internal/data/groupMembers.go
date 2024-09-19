package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/manuelam2003/triclone/internal/validator"
)

type GroupMember struct {
	ID       int64      `json:"id"`
	GroupID  int64      `json:"group_id"`
	UserID   int64      `json:"user_id"`
	JoinedAt time.Time  `json:"joined_at"`
	IsActive bool       `json:"is_active"`
	LeftAt   *time.Time `json:"left_at"`
}

func ValidateGroupMember(v *validator.Validator, groupMember *GroupMember) {
	v.Check(groupMember.GroupID > 0, "group_id", "must be non negative")
	v.Check(groupMember.UserID > 0, "user_id", "must be non negative")
}

type GroupMemberModel struct {
	DB *sql.DB
}

func (m GroupMemberModel) Insert(groupID, userID int64) error {
	query := `
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, groupID, userID)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "group_members_group_id_user_id_key"`:
			return ErrDuplicateEntry
		case err.Error() == `pq: insert or update on table "group_members" violates foreign key constraint "group_members_group_id_fkey"`:
			return ErrForeignKeyViolation
		default:
			return err
		}
	}
	return nil
}

func (m GroupMemberModel) SoftDelete(groupID, userID int64) error {
	query := `
		UPDATE group_members
		SET is_active = FALSE, left_at = $1
		WHERE group_id = $2 AND user_id = $3 AND is_active = TRUE`

	leftAt := time.Now()

	result, err := m.DB.Exec(query, leftAt, groupID, userID)
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

func (m GroupMemberModel) UserBelongsToGroup(userID, groupID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM group_members 
			WHERE user_id = $1 AND group_id = $2 AND is_active = true
		)`

	var exists bool

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID, groupID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m GroupMemberModel) CheckIfUserWasInGroup(groupID, userID int64) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM group_members 
		WHERE group_id = $1 AND user_id = $2 AND is_active = false`

	var count int
	err := m.DB.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m GroupMemberModel) ReinstateMember(groupID, userID int64) error {
	query := `
		UPDATE group_members
		SET is_active = true, left_at = NULL
		WHERE group_id = $1 AND user_id = $2 AND is_active = false`

	_, err := m.DB.Exec(query, groupID, userID)
	return err
}
