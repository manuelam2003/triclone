package data

import "database/sql"

type Balance struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

type BalanceModel struct {
	DB *sql.DB
}

func (m BalanceModel) CalculateGroupBalances(groupID int64) ([]Balance, error) {
	// Query to calculate the total amount each user owes based on expenses
	expensesQuery := `
        SELECT p.user_id, SUM(p.amount_owed) AS total_owed
        FROM expense_participants p
        INNER JOIN expenses e ON p.expense_id = e.id
        WHERE e.group_id = $1
        GROUP BY p.user_id
    `

	// Query to calculate the total settled amounts for each user
	settlementsQuery := `
        SELECT s.payer_id, s.payee_id, SUM(s.amount) AS settled_amount
        FROM settlements s
        WHERE s.group_id = $1
        GROUP BY s.payer_id, s.payee_id
    `

	expenseMap := make(map[int64]float64)
	var balances []Balance

	// Query expenses and build the map of user balances based on what they owe
	rows, err := m.DB.Query(expensesQuery, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int64
		var totalOwed float64
		if err := rows.Scan(&userID, &totalOwed); err != nil {
			return nil, err
		}
		expenseMap[userID] = totalOwed
	}

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
