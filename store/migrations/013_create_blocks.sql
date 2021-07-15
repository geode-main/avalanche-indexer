-- +goose Up
CREATE TYPE block_type AS ENUM (
  'proposal',
  'standard',
  'atomic',
  'commit',
  'abort',
  'evm'
);

CREATE TABLE blocks (
  id        TEXT NOT NULL PRIMARY KEY,
  parent    TEXT NOT NULL,
  chain     TEXT NOT NULL,
  height    INTEGER NOT NULL,
  type      block_type NOT NULL,
  processed BOOLEAN DEFAULT FALSE,
  timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_blocks_chain_height ON blocks (chain, height);
CREATE INDEX idx_blocks_type         ON blocks(type);
CREATE INDEX idx_blocks_time         ON blocks(timestamp);

-- +goose Down
DROP TABLE blocks;
DROP TYPE block_type;