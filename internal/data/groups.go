package data

import (
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
