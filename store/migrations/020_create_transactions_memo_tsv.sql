-- +goose Up
ALTER TABLE transactions ADD COLUMN memo_tsv TSVECTOR;

CREATE INDEX CONCURRENTLY idx_transactions_memo_tsv ON transactions USING GIN (memo_tsv);

UPDATE transactions
SET memo_tsv = to_tsvector(memo_text)
WHERE
  memo_text IS NOT NULL
  AND LENGTH(memo_text) > 1;

-- +goose Down
ALTER TABLE transactions DROP COLUMN memo_tsv;