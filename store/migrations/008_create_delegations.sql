-- +goose Up
CREATE TABLE delegations (
  id                       SERIAL PRIMARY KEY,
  reference_id             VARCHAR(64) NOT NULL,
  node_id                  VARCHAR(64) NOT NULL,
  reward_address           VARCHAR(64) NOT NULL,
  stake_amount             BIGINT,
  potential_reward         BIGINT,
  active                   BOOLEAN,
  active_start_time        TIMESTAMP WITH TIME ZONE,
  active_end_time          TIMESTAMP WITH TIME ZONE,
  active_progress_percent  DECIMAL,
  first_height             INTEGER,
  last_height              INTEGER,
  created_at               TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at               TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_delegations_reference
  ON delegations(reference_id);

CREATE INDEX idx_delegations_node_id
  ON delegations(node_id);

CREATE INDEX idx_delegations_address
  ON delegations(reward_address);

-- +goose Down
DROP TABLE delegations;
