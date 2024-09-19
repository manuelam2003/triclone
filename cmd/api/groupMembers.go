package main

import (
	"errors"
	"net/http"

	"github.com/manuelam2003/triclone/internal/data"
	"github.com/manuelam2003/triclone/internal/validator"
)

func (app *application) addGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	v := validator.New()

	err = app.models.GroupMembers.Insert(groupID, currentUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEntry):
			v.AddError("member", "this user is already a member of this group")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrForeignKeyViolation):
			v.AddError("group_id", "a group with this group_id does not exist")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
}

func (app *application) listGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	// groupID, err := app.readIDParam(r, "group_id")
	// if err != nil {
	// 	app.notFoundResponse(w, r)
	// 	return
	// }
}

func (app *application) removeGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	// groupID, err := app.readIDParam(r, "group_id")
	// if err != nil {
	// 	app.notFoundResponse(w, r)
	// 	return
	// }
}
