-- +goose Up
CREATE TABLE evm_receipts (
  id               SERIAL PRIMARY KEY,
  block_height     INTEGER NOT NULL,
  contract_address TEXT
);

CREATE TABLE evm_logs (
  id         SERIAL PRIMARY KEY,
  receipt_id INTEGER NOT NULL,
  idx        INTEGER,
  address    TEXT,
  tx_idx     INTEGER,
  removed    BOOLEAN,
  topics     TEXT[],
  data       TEXT
);

-- +goose Down
DROP TABLE IF EXISTS evm_receipts;
DROP TABLE IF EXISTS evm_logs;