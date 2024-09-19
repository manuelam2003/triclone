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

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully added to group"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	group, err := app.models.Groups.Get(groupID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	userID, err := app.readIDParam(r, "user_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	// ! cuidao si es nil el pointer
	if currentUser.ID != userID && currentUser.ID != *group.CreatedBy {
		app.invalidUserResponse(w, r)
		return
	}

	err = app.models.GroupMembers.SoftDelete(groupID, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully removed from group"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) reinstateGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	userID, err := app.readIDParam(r, "user_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	exists, err := app.models.GroupMembers.CheckIfUserWasInGroup(groupID, userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !exists {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.GroupMembers.ReinstateMember(groupID, userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Member reinstated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
