package main

import (
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

	// Check if the expense belongs to the specified group
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

func (app *application) addExpenseParticipantHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateExpenseParticipantHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteExpenseParticipantHandler(w http.ResponseWriter, r *http.Request) {

}
