-- +goose Up
CREATE TYPE transaction_status AS ENUM (
  'accepted',
  'rejected',
  'reverted'
);

CREATE TYPE transaction_type AS ENUM (
  'p_create_chain',
	'p_create_subnet',
	'p_add_subnet_validator',
	'p_advance_time',
	'p_reward_validator',
	'p_add_validator',
	'p_add_delegator',
	'p_import',
	'p_export',
  'x_base',
  'x_import',
  'x_export',
  'x_create_asset',
  'x_operation',
  'c_atomic_export',
  'c_atomic_import',
  'c_evm'
);

CREATE TABLE transactions (
  id                TEXT NOT NULL PRIMARY KEY,
  type              transaction_type NOT NULL,
  status            transaction_status NOT NULL,
  block             TEXT,
  block_height      INTEGER,
  chain             TEXT NOT NULL,
  memo              TEXT,
  memo_text         TEXT,
  nonce             INTEGER,
  fee               BIGINT,
  source_chain      TEXT,
  destination_chain TEXT,
  reference_tx_id   TEXT,
  metadata          JSONB,
  timestamp         TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_transactions_type         ON transactions(type);
CREATE INDEX idx_transactions_status       ON transactions(status);
CREATE INDEX idx_transactions_chain        ON transactions(chain);
CREATE INDEX idx_transactions_time         ON transactions(timestamp);
CREATE INDEX idx_transactions_block_hash   ON transactions(block);
CREATE INDEX idx_transactions_block_height ON transactions(block_height);
CREATE INDEX idx_transactions_memo_text    ON transactions USING GIN (to_tsvector('english', memo_text));

-- +goose Down
DROP TABLE transactions;

DROP TYPE transaction_type;
DROP TYPE transaction_status;