package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Groups       GroupModel
	Users        UserModel
	Tokens       TokenModel
	GroupMembers GroupMemberModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Groups:       GroupModel{DB: db},
		Users:        UserModel{DB: db},
		Tokens:       TokenModel{DB: db},
		GroupMembers: GroupMemberModel{DB: db},
	}
}
