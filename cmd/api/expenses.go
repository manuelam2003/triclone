package main

import (
	"fmt"
	"net/http"

	"github.com/manuelam2003/triclone/internal/data"
	"github.com/manuelam2003/triclone/internal/validator"
)

func (app *application) listGroupExpensesHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) showGroupExpenseHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createGroupExpenseHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// TODO only current user can make an expense
	var input struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.models.GroupMembers.UserBelongsToGroup(currentUser.ID, groupID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !isMember {
		app.forbiddenResponse(w, r)
		return
	}

	expense := &data.Expense{
		GroupID:     groupID,
		Amount:      input.Amount,
		Description: input.Description,
		PaidBy:      &currentUser.ID,
	}

	v := validator.New()

	if data.ValidateExpense(v, expense); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Expenses.Insert(expense)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/groups/%d/expenses/%d", groupID, expense.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"expense": expense}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateGroupExpenseHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteGroupExpenseHandler(w http.ResponseWriter, r *http.Request) {

}
