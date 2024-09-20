package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/manuelam2003/triclone/internal/data"
	"github.com/manuelam2003/triclone/internal/validator"
)

func (app *application) listExpenseParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	ids, err := app.extractIDsFromRequest(r, "group_id", "expense_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, ids["group_id"])
	if err != nil || !isMember {
		return
	}

	if _, err := app.checkExpenseInGroup(w, r, ids["expense_id"], ids["group_id"]); err != nil {
		return
	}

	var input struct{ data.Filters }
	v := validator.New()

	qs := r.URL.Query()
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "amount_owed", "user_id", "-id", "-amount_owed", "-user_id"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	participants, metadata, err := app.models.ExpensesParticipants.GetAllForGroupAndExpense(ids["group_id"], ids["expense_id"], input.Filters)
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
	ids, err := app.extractIDsFromRequest(r, "group_id", "expense_id")
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

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, ids["group_id"])
	if err != nil || !isMember {
		return
	}

	if _, err := app.checkExpenseInGroup(w, r, ids["expense_id"], ids["group_id"]); err != nil {
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
	var invalidParticipants []int64
	var membershipErrors []int64
	v := validator.New()

	for _, participant := range participants {

		isMember, err := app.checkUserMembership(w, r, participant.UserID, ids["group_id"])
		if err != nil || !isMember {
			membershipErrors = append(membershipErrors, participant.UserID)
			continue
		}

		newParticipant := &data.ExpenseParticipant{
			ExpenseID:  expenseID,
			UserID:     participant.UserID,
			AmountOwed: participant.AmountOwed,
		}

		if data.ValidateParticipant(v, newParticipant); !v.Valid() {
			invalidParticipants = append(invalidParticipants, participant.UserID)
			continue
		}

		err = app.models.ExpensesParticipants.Insert(newParticipant)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		newRecords++
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "expense participants added succesfully",
		"newRecords":          newRecords,
		"invalidParticipants": invalidParticipants,
		"membershipErrors":    membershipErrors,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateExpenseParticipantHandler(w http.ResponseWriter, r *http.Request) {
	ids, err := app.extractIDsFromRequest(r, "group_id", "expense_id", "participant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, ids["group_id"])
	if err != nil || !isMember {
		return
	}

	if _, err := app.checkExpenseInGroup(w, r, ids["expense_id"], ids["group_id"]); err != nil {
		return
	}

	// Fetch the participant based on the participant ID
	participant, err := app.models.ExpensesParticipants.Get(ids["participant_id"])
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if participant.ExpenseID != ids["expense_id"] {
		app.notFoundResponse(w, r)
		return
	}

	isMember, err = app.checkUserMembership(w, r, participant.UserID, ids["group_id"])
	if err != nil || !isMember {
		return
	}

	var input struct {
		AmountOwed float64 `json:"amount_owed"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	participant.AmountOwed = input.AmountOwed

	err = app.models.ExpensesParticipants.Update(participant)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"participant": participant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteExpenseParticipantHandler(w http.ResponseWriter, r *http.Request) {
	ids, err := app.extractIDsFromRequest(r, "group_id", "expense_id", "participant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, ids["group_id"])
	if err != nil || !isMember {
		return
	}

	if _, err := app.checkExpenseInGroup(w, r, ids["expense_id"], ids["group_id"]); err != nil {
		return
	}

	participant, err := app.models.ExpensesParticipants.Get(ids["participant_id"])
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if participant.ExpenseID != ids["expense_id"] {
		app.notFoundResponse(w, r)
		return
	}

	isMember, err = app.checkUserMembership(w, r, participant.UserID, ids["group_id"])
	if err != nil || !isMember {
		return
	}

	err = app.models.ExpensesParticipants.Delete(ids["participant_id"])
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "participant successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Helper function to check if an expense belongs to the group
func (app *application) checkExpenseInGroup(w http.ResponseWriter, r *http.Request, expenseID, groupID int64) (bool, error) {
	belongsToGroup, err := app.models.Expenses.CheckExpenseBelongsToGroup(expenseID, groupID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return false, err
	}
	if !belongsToGroup {
		app.notFoundResponse(w, r)
		return false, nil
	}
	return true, nil
}
