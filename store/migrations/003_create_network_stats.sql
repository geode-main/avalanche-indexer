-- +goose Up
CREATE TABLE network_stats (
  id                  SERIAL PRIMARY KEY,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  bucket              VARCHAR(16),
  height_change       INTEGER,
  peers               INTEGER,
  blockchains         INTEGER,
  active_validators   INTEGER,
  pending_validators  INTEGER,
  active_delegations  INTEGER,
  pending_delegations INTEGER,
  uptime              DECIMAL
);

CREATE UNIQUE INDEX idx_network_stats_timebucket
  ON network_stats(time, bucket);

-- +goose Down
DROP TABLE network_stats;
