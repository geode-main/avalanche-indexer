-- +goose Up
CREATE TABLE network_metrics (
  id                        SERIAL PRIMARY KEY,
  time                      TIMESTAMP WITH TIME ZONE,
  height                    INTEGER,
  peers_count               INTEGER,
  blockchains_count         INTEGER,
  active_validators_count   INTEGER,
  pending_validators_count  INTEGER,
  active_delegations_count  INTEGER,
  pending_delegations_count INTEGER,
  min_validator_stake       BIGINT,
  min_delegation_stake      BIGINT,
  tx_fee                    INT,
  creation_tx_fee           INT,
  uptime                    DECIMAL,
  delegation_fee            DECIMAL
);

CREATE INDEX idx_network_metrics_time
  ON network_metrics(time);

-- +goose Down
DROP TABLE network_metrics;
