-- +goose Up
ALTER TABLE transaction_outputs ALTER COLUMN amount TYPE DECIMAL(65, 0);

-- +goose Up
SELECT 1;