-- +goose Up
CREATE TABLE rewards (
  id             TEXT NOT NULL PRIMARY KEY,
  transaction_id TEXT NOT NULL,
  rewarded       BOOLEAN NOT NULL,
  rewarded_at    TIMESTAMP WITH TIME ZONE,
  processed_at   TIMESTAMP WITH TIME ZONE
);

CREATE TABLE rewards_owners (
  id        TEXT NOT NULL PRIMARY KEY,
  locktime  INTEGER NOT NULL,
  threshold INTEGER NOT NULL
);

CREATE TABLE rewards_owner_addresses (
  id      TEXT NOT NULL,
  address TEXT NOT NULL,
  index   INTEGER NOT NULL
);

CREATE TABLE rewards_owner_outputs (
  id             TEXT NOT NULL PRIMARY KEY,
  transaction_id TEXT NOT NULL,
  index          INTEGER NOT NULL
);

CREATE INDEX idx_rewards_owner_addresses_id
  ON rewards_owner_addresses(id);

CREATE UNIQUE INDEX idx_rewards_owner_outputs_id
  ON rewards_owner_outputs(id);

CREATE INDEX idx_rewards_owner_outputs_tx_id
  ON rewards_owner_outputs(transaction_id);

-- +goose Down
DROP TABLE rewards;
DROP TABLE rewards_owner_outputs;
DROP TABLE rewards_owner_addresses;
DROP TABLE rewards_owners;
