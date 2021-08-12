-- +goose Up
DROP TABLE IF EXISTS evm_receipts;
DROP TABLE IF EXISTS evm_logs;

CREATE TABLE evm_receipts (
  id               TEXT PRIMARY KEY,
  type             INTEGER NOT NULL,
  status           INTEGER NOT NULL,
  contract_address TEXT,
  logs             JSONB
);

CREATE TABLE evm_traces (
  id        TEXT PRIMARY KEY,
  data      TEXT NOT NULL,
  timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE evm_receipts;
DROP TABLE evm_traces;
