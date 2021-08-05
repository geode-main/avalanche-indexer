-- +goose Up
DROP INDEX idx_transactions_memo_text;

-- +goose Down
CREATE INDEX idx_transactions_memo_text ON transactions USING GIN (to_tsvector('english', memo_text));