package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/manuelam2003/triclone/internal/data"
	"github.com/manuelam2003/triclone/internal/validator"
)

func (app *application) listGroupsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string
		CreatedBy int64
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.CreatedBy = int64(app.readInt(qs, "created_by", 0, v))

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "created_by", "-id", "-name", "-created_by"}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	groups, metadata, err := app.models.Groups.GetAll(input.Name, input.CreatedBy, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "groups": groups}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createGroupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	currentUser := app.contextGetUser(r)

	// ! cuidao si es nil el pointer
	group := &data.Group{
		Name:      input.Name,
		CreatedBy: &currentUser.ID,
	}

	v := validator.New()

	if data.ValidateGroup(v, group); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Groups.Insert(group)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/groups/%d", group.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"group": group}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	group, err := app.models.Groups.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"group": group}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	group, err := app.models.Groups.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	currentUser := app.contextGetUser(r)

	// ! cuidao si es nil el pointer
	if currentUser.ID != *group.CreatedBy {
		app.invalidUserResponse(w, r)
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Groups.Update(group)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"group": group}, nil)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
}

func (app *application) deleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	group, err := app.models.Groups.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	currentUser := app.contextGetUser(r)

	// ! cuidao si es nil el pointer
	if currentUser.ID != *group.CreatedBy {
		app.invalidUserResponse(w, r)
		return
	}

	err = app.models.Groups.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "group succesfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
