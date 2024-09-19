-- 1. Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated boolean NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Groups Table
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Group Members Table
CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,  
    left_at TIMESTAMP,  
    UNIQUE (group_id, user_id)
);

-- 4. Expenses Table
CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    amount NUMERIC(10, 2) NOT NULL,
    description VARCHAR(255),
    paid_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Expense Participants Table
CREATE TABLE expense_participants (
    id SERIAL PRIMARY KEY,
    expense_id INT REFERENCES expenses(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    amount_owed NUMERIC(10, 2) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (expense_id, user_id)
);

-- 6. Settlements Table
CREATE TABLE settlements (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    payer_id INT REFERENCES users(id) ON DELETE SET NULL,
    payee_id INT REFERENCES users(id) ON DELETE SET NULL,
    amount NUMERIC(10, 2) NOT NULL,
    settled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (payer_id != payee_id)
);

CREATE TABLE tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);

CREATE INDEX idx_group_id ON group_members(group_id);
CREATE INDEX idx_expense_id ON expense_participants(expense_id);
CREATE INDEX idx_group_expenses ON expenses(group_id);
CREATE INDEX idx_group_settlements ON settlements(group_id);