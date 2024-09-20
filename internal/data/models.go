package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrEditConflict        = errors.New("edit conflict")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
)

type Models struct {
	Groups               GroupModel
	Users                UserModel
	Tokens               TokenModel
	GroupMembers         GroupMemberModel
	Expenses             ExpenseModel
	ExpensesParticipants ExpenseParticipantModel
	Settlements          SettlementModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Groups:               GroupModel{DB: db},
		Users:                UserModel{DB: db},
		Tokens:               TokenModel{DB: db},
		GroupMembers:         GroupMemberModel{DB: db},
		Expenses:             ExpenseModel{DB: db},
		ExpensesParticipants: ExpenseParticipantModel{DB: db},
		Settlements:          SettlementModel{DB: db},
	}
}
