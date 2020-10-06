-- +goose Up
CREATE TABLE addresses (
  id                   SERIAL PRIMARY KEY,
  value                TEXT NOT NULL,
  balance              BIGINT,
  unlocked_balance     BIGINT,
  locked_not_stakeable BIGINT,
  created_at           TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at           TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_addresses_value
  ON addresses (value);

-- +goose Down
DROP TABLE addresses;
