package main

import (
	"encoding/json"
	"net/http"

	"github.com/manuelam2003/triclone/internal/data"
	"github.com/manuelam2003/triclone/internal/validator"
)

func (app *application) listExpenseParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	expenseID, err := app.readIDParam(r, "expense_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, groupID)
	if err != nil || !isMember {
		return
	}

	belongsToGroup, err := app.models.Expenses.CheckExpenseBelongsToGroup(expenseID, groupID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !belongsToGroup {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "amount_owed", "user_id", "-id", "-amount_owed", "-user_id"}

	participants, metadata, err := app.models.ExpensesParticipants.GetAllForGroupAndExpense(groupID, expenseID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"participants": participants, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addExpenseParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	expenseID, err := app.readIDParam(r, "expense_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, groupID)
	if err != nil || !isMember {
		return
	}

	belongsToGroup, err := app.models.Expenses.CheckExpenseBelongsToGroup(expenseID, groupID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !belongsToGroup {
		app.notFoundResponse(w, r)
		return
	}

	var participants []struct {
		UserID     int64   `json:"user_id"`
		AmountOwed float64 `json:"amount_owed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&participants); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	newRecords := 0

	v := validator.New()

	for _, participant := range participants {

		isMember, err := app.checkUserMembership(w, r, participant.UserID, groupID)
		if err != nil || !isMember {
			continue
		}

		newParticipant := &data.ExpenseParticipant{
			ExpenseID:  expenseID,
			UserID:     participant.UserID,
			AmountOwed: participant.AmountOwed,
		}

		if data.ValidateParticipant(v, newParticipant); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.models.ExpensesParticipants.Insert(newParticipant)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		newRecords++
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "expense participants added succesfully", "newRecords": newRecords}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateExpenseParticipantHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteExpenseParticipantHandler(w http.ResponseWriter, r *http.Request) {

}
