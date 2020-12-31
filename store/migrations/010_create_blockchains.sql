-- +goose Up
CREATE TABLE blockchains (
  id        TEXT PRIMARY KEY,
  name      TEXT NOT NULL,
  subnet_id TEXT NOT NULL,
  vm_id     TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS blockchains;