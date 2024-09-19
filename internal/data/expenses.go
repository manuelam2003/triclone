package data

import "time"

type Expense struct {
	ID          int64     `json:"id"`          // SERIAL (int in Go)
	GroupID     int64     `json:"group_id"`    // INT (references groups(id))
	Amount      float64   `json:"amount"`      // NUMERIC(10, 2) (float64 in Go)
	Description string    `json:"description"` // VARCHAR(255) (string in Go)
	PaidBy      *int64    `json:"paid_by"`     // INT (references users(id)) - can be NULL, so use a pointer
	CreatedAt   time.Time `json:"created_at"`  // TIMESTAMP (time.Time in Go)
	UpdatedAt   time.Time `json:"updated_at"`  // TIMESTAMP (time.Time in Go)
}

// -- 4. Expenses Table
// CREATE TABLE expenses (
//     id SERIAL PRIMARY KEY,
//     group_id INT REFERENCES groups(id) ON DELETE CASCADE,
//     amount NUMERIC(10, 2) NOT NULL,
//     description VARCHAR(255),
//     paid_by INT REFERENCES users(id) ON DELETE SET NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
