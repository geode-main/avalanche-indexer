-- +goose Up
CREATE TABLE validator_sequences (
  id                       SERIAL PRIMARY KEY,
  time                     TIMESTAMP WITH TIME ZONE,
  height                   INTEGER NOT NULL,
  node_id                  VARCHAR(64) NOT NULL,
  stake_amount             BIGINT,
  stake_percent            DECIMAL,
  potential_reward         BIGINT,
  reward_address           VARCHAR(64),
  active                   BOOLEAN,
  active_start_time        TIMESTAMP WITH TIME ZONE,
  active_end_time          TIMESTAMP WITH TIME ZONE,
  active_progress_percent  DECIMAL,
  delegations_count        INTEGER,
  delegations_percent      DECIMAL,
  delegated_amount         BIGINT,
  delegated_amount_percent DECIMAL,
  delegation_fee           DECIMAL,
  uptime                   DECIMAL
);

CREATE INDEX idx_validator_seq_height
  ON validator_sequences(height);

CREATE index idx_validator_seq_time
  ON validator_sequences(time);

-- +goose Down
DROP TABLE validator_sequences;
