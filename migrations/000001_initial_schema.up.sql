-- 1. Users Table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Groups Table
CREATE TABLE groups (
    group_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_by INT REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Group Members Table
CREATE TABLE group_members (
    group_member_id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(group_id) ON DELETE CASCADE,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (group_id, user_id)
);

-- 4. Expenses Table
CREATE TABLE expenses (
    expense_id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(group_id) ON DELETE CASCADE,
    amount NUMERIC(10, 2) NOT NULL,
    description VARCHAR(255),
    paid_by INT REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Expense Participants Table
CREATE TABLE expense_participants (
    expense_participant_id SERIAL PRIMARY KEY,
    expense_id INT REFERENCES expenses(expense_id) ON DELETE CASCADE,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    amount_owed NUMERIC(10, 2) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (expense_id, user_id)
);

-- 6. Settlements Table
CREATE TABLE settlements (
    settlement_id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(group_id) ON DELETE CASCADE,
    payer_id INT REFERENCES users(user_id) ON DELETE SET NULL,
    payee_id INT REFERENCES users(user_id) ON DELETE SET NULL,
    amount NUMERIC(10, 2) NOT NULL,
    settled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (payer_id != payee_id)
);

CREATE INDEX idx_group_id ON group_members(group_id);
CREATE INDEX idx_expense_id ON expense_participants(expense_id);
CREATE INDEX idx_group_expenses ON expenses(group_id);
CREATE INDEX idx_group_settlements ON settlements(group_id);