-- +goose Up
CREATE TABLE chains (
  id         SERIAL PRIMARY KEY,
  chain_id   TEXT NOT NULL,
  vm         TEXT,
  name       TEXT,
  network    INTEGER,
  subnet     TEXT
);

CREATE UNIQUE index idx_chains_id ON chains(chain_id);

CREATE TABLE sync_statuses (
  id         TEXT PRIMARY KEY,
  index_id   INTEGER,
  index_time TIMESTAMP WITH TIME ZONE,
  tip_id     INTEGER,
  tip_time   TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE sync_statuses;
DROP TABLE chains;