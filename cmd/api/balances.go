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

// func (app *application) calculateBalancesHandler(w http.ResponseWriter, r *http.Request) {
// 	groupID, err := app.readIDParam(r, "group_id")
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	// Fetch all expenses for the group
// 	expenses, _, err := app.models.Expenses.GetAll(groupID, "", 0, data.Filters{})
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	// Initialize a map to hold balances
// 	balances := make(map[int64]float64)

// 	for _, expense := range expenses {
// 		// Get participants for the expense
// 		participants, err := app.models.ExpensesParticipants.GetAllForExpense(expense.ID)
// 		if err != nil {
// 			app.serverErrorResponse(w, r, err)
// 			return
// 		}

// 		// Calculate balance based on the paid_by
// 		for _, participant := range participants {
// 			if participant.UserID == expense.PaidBy {
// 				// If the participant is the one who paid, subtract the total expense amount
// 				balances[participant.UserID] -= expense.Amount
// 			} else {
// 				// Otherwise, add the amount owed
// 				balances[participant.UserID] += participant.AmountOwed
// 			}
// 		}
// 	}

// 	// Convert balances map to a slice for JSON response
// 	balanceSlice := make([]struct {
// 		UserID  int64   `json:"user_id"`
// 		Balance float64 `json:"balance"`
// 	}, 0)

// 	for userID, balance := range balances {
// 		balanceSlice = append(balanceSlice, struct {
// 			UserID  int64   `json:"user_id"`
// 			Balance float64 `json:"balance"`
// 		}{UserID: userID, Balance: balance})
// 	}

// 	// Write the response
// 	err = app.writeJSON(w, http.StatusOK, envelope{"balances": balanceSlice}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }
