package data

import (
	"database/sql"
	"fmt"
)

type Balance struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

type BalanceModel struct {
	DB *sql.DB
}

func (m BalanceModel) CalculateGroupBalances(groupID int64) ([]Balance, error) {
	expensesQuery := `
        SELECT e.id, e.paid_by, p.user_id, SUM(p.amount_owed) AS total_owed, e.amount AS total_paid
        FROM expense_participants p
        INNER JOIN expenses e ON p.expense_id = e.id
        WHERE e.group_id = $1
        GROUP BY e.id, e.paid_by, p.user_id, e.amount
    `

	settlementsQuery := `
        SELECT s.payer_id, s.payee_id, SUM(s.amount) AS settled_amount
        FROM settlements s
        WHERE s.group_id = $1
        GROUP BY s.payer_id, s.payee_id
    `

	processedExpenses := make(map[int64]struct{}) // Set to track processed expense IDs
	expenseMap := make(map[int64]float64)
	var balances []Balance

	// Query expenses and build the map of user balances based on what they owe
	rows, err := m.DB.Query(expensesQuery, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var expenseID, paidByUserID, participantUserID int64
		var totalOwed, totalPaid float64
		if err := rows.Scan(&expenseID, &paidByUserID, &participantUserID, &totalOwed, &totalPaid); err != nil {
			return nil, err
		}

		fmt.Printf("Expense: %d, Paid by: %d, User id: %d, Total owed: %.2f, Total amount: %.2f\n", expenseID, paidByUserID, participantUserID, totalOwed, totalPaid)

		// Increase the balance of the participants by what they owe
		expenseMap[participantUserID] += totalOwed

		if _, exists := processedExpenses[expenseID]; !exists {
			// Reduce the balance of the person who paid by the full amount of the expense
			processedExpenses[expenseID] = struct{}{}
			expenseMap[paidByUserID] -= totalPaid
			continue
		}

	}
	fmt.Println(expenseMap)

	// Query settlements and adjust balances
	rows, err = m.DB.Query(settlementsQuery, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payerID, payeeID int64
		var settledAmount float64
		if err := rows.Scan(&payerID, &payeeID, &settledAmount); err != nil {
			return nil, err
		}

		// Deduct settled amount from payer's debt
		expenseMap[payerID] -= settledAmount
		// Add the settled amount to payee's debt (or reduce how much they are owed)
		expenseMap[payeeID] += settledAmount
	}

	// Convert the expenseMap to the slice of balances
	for userID, balance := range expenseMap {
		balances = append(balances, Balance{
			UserID:  userID,
			Balance: balance,
		})
	}

	return balances, nil
}
