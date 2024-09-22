package main

import (
	"net/http"
)

func (app *application) groupBalanceHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, groupID)
	if err != nil || !isMember {
		app.forbiddenResponse(w, r)
		return
	}

	balances, err := app.models.Balances.CalculateGroupBalances(groupID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"balances": balances}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
