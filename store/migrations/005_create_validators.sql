-- +goose Up
CREATE TABLE validators (
  id                       BIGSERIAL PRIMARY KEY,
  node_id                  VARCHAR(64) NOT NULL,
  stake_amount             BIGINT,
  stake_percent            DECIMAL,
  potential_reward         BIGINT,
  reward_address           VARCHAR(64),
  active                   BOOLEAN,
  active_start_time        TIMESTAMP WITH TIME ZONE,
  active_end_time          TIMESTAMP WITH TIME ZONE,
  active_progress_percent  DECIMAL,
  uptime                   DECIMAL,
  delegations_count        INTEGER,
  delegations_percent      DECIMAL,
  delegated_amount         BIGINT,
  delegated_amount_percent DECIMAL,
  delegation_fee           DECIMAL,
  capacity                 BIGINT,
  capacity_percent         DECIMAL,
  first_height             INTEGER,
  last_height              INTEGER,
  created_at               TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at               TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_validators_node_id
  ON validators(node_id);

CREATE INDEX idx_validators_stake_amount
  ON validators(stake_amount);

-- +goose Down
DROP TABLE validators;
