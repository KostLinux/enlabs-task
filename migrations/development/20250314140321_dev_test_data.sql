-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- This migration only runs in development environment
INSERT INTO users (id, username) VALUES 
    (1, 'user1'),
    (2, 'user2'),
    (3, 'user3');

INSERT INTO balances (user_id, amount) VALUES 
    (1, 1000.00),
    (2, 1000.00),
    (3, 1000.00);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DELETE FROM users WHERE id IN (1, 2, 3);
-- The balance entries will be automatically deleted due to ON DELETE CASCADE