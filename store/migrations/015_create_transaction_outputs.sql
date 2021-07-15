-- +goose Up
CREATE TYPE output_type AS ENUM (
  'transfer',
  'stakeable_lock',
  'mint',
  'nft_mint',
  'nft_transfer'
);

CREATE TABLE transaction_inputs (
  id    TEXT NOT NULL PRIMARY KEY,
  tx_id TEXT NOT NULL
);

CREATE TABLE transaction_outputs (
  id          TEXT NOT NULL PRIMARY KEY,
  tx_id       TEXT NOT NULL,
  chain       TEXT NOT NULL,
  asset       TEXT NOT NULL,
  type        output_type NOT NULL,
  index       INTEGER NOT NULL,
  locktime    BIGINT NOT NULL,
  threshold   INTEGER NOT NULL,
  amount      BIGINT,
  group_id    INTEGER,
  stake       BOOLEAN NOT NULL,
  reward      BOOLEAN NOT NULL,
  spent       BOOLEAN NOT NULL,
  spent_tx_id TEXT,
  payload     TEXT,
  addresses   TEXT[]
);

CREATE INDEX idx_transaction_outputs_tx        ON transaction_outputs(tx_id);
CREATE INDEX idx_transaction_outputs_spent_tx  ON transaction_outputs(spent_tx_id);
CREATE INDEX idx_transaction_outputs_type      ON transaction_outputs(type);
CREATE INDEX idx_transaction_outputs_addresses ON transaction_outputs USING GIN(addresses);
CREATE INDEX idx_transaction_outputs_asset     ON transaction_outputs(asset);

-- +goose Down
DROP TABLE transaction_inputs;
DROP TABLE transaction_outputs;

DROP TYPE output_type;