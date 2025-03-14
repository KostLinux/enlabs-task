-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE balance (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    amount DECIMAL(19, 2) NOT NULL DEFAULT 0.00,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_balance UNIQUE (user_id)
);

CREATE INDEX idx_balance_user_id ON balance(user_id);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    state VARCHAR(10) NOT NULL,
    amount DECIMAL(19, 2) NOT NULL,
    source_type VARCHAR(20) NOT NULL,
    previous_balance DECIMAL(19, 2) NOT NULL,
    new_balance DECIMAL(19, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_transaction_id UNIQUE (transaction_id)
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_transaction_id ON transactions(transaction_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS balance;
DROP TABLE IF EXISTS users;