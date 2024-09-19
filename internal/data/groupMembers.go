package data

import (
	"database/sql"
	"time"

	"github.com/manuelam2003/triclone/internal/validator"
)

type GroupMember struct {
	ID       int64     `json:"id"`
	GroupID  int64     `json:"group_id"`
	UserID   int64     `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
	IsActive bool      `json:"is_active"`
	LeftAt   time.Time `json:"left_at"`
}

func ValidateGroupMember(v *validator.Validator, groupMember *GroupMember) {
	v.Check(groupMember.ID > 0, "id", "must be non negative")
	v.Check(groupMember.GroupID > 0, "group_id", "must be non negative")
	v.Check(groupMember.UserID > 0, "user_id", "must be non negative")
}

type GroupMemberModel struct {
	DB *sql.DB
}

func ()