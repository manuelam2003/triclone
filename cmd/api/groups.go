package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/manuelam2003/triclone/internal/data"
	"github.com/manuelam2003/triclone/internal/validator"
)

func (app *application) createGroupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		CreatedBy int64  `json:"created_by"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	group := &data.Group{
		Name:      input.Name,
		CreatedBy: input.CreatedBy,
	}

	v := validator.New()

	if data.ValidateGroup(v, group); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showGroupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	group := data.Group{
		ID:        id,
		Name:      "Grupo 1",
		CreatedBy: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"group": group}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
