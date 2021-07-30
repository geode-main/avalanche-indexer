-- +goose Up
CREATE TABLE events (
  id           TEXT NOT NULL PRIMARY KEY,
  scope        TEXT NOT NULL,
  type         TEXT NOT NULL,
  chain        TEXT NOT NULL,
  block_hash   TEXT,
  block_height INTEGER,
  tx_hash      TEXT,
  item_id      TEXT,
  item_type    TEXT,
  data         JSONB,
  timestamp    TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_events_scope ON events(scope);
CREATE INDEX idx_events_chain ON events(chain);
CREATE INDEX idx_events_type  ON events(type);
CREATE INDEX idx_events_item  ON events(item_id, item_type);
CREATE INDEX idx_events_block ON events(block_height);
CREATE INDEX idx_events_time  ON events(timestamp);

-- +goose Down
DROP TABLE events;