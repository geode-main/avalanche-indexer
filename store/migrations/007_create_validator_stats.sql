-- +goose Up
CREATE TABLE validator_stats (
  id                       BIGSERIAL PRIMARY KEY,
  time                     TIMESTAMP WITH TIME ZONE NOT NULL,
  bucket                   VARCHAR(16) NOT NULL,
  node_id                  VARCHAR(64) NOT NULL,
  uptime_min               DECIMAL,
  uptime_max               DECIMAL,
  uptime_avg               DECIMAL,
  stake_amount             BIGINT,
  stake_percent            DECIMAL,
  delegations_count        INTEGER,
  delegations_percent      DECIMAL,
  delegated_amount         BIGINT,
  delegated_amount_percent DECIMAL
);

CREATE INDEX idx_validator_stats_time
  ON validator_stats (time);

-- +goose Down
DROP TABLE validator_stats;
