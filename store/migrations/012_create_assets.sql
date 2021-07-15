-- +goose Up
CREATE TYPE asset_type AS ENUM (
  'fixed_cap',
  'variable_cap',
  'nft'
);

CREATE TABLE assets (
  id           SERIAL PRIMARY KEY,
  asset_id     TEXT NOT NULL,
  type         asset_type NOT NULL,
  name         TEXT,
  symbol       TEXT,
  denomination INTEGER
);

CREATE UNIQUE INDEX idx_assets_id ON assets(asset_id);
CREATE INDEX idx_assets_type ON assets(type);

-- +goose Down
DROP TABLE assets;
DROP TYPE asset_type;