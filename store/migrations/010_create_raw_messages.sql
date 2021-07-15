-- +goose Up
CREATE TABLE raw_message_topics (
  id    SERIAL PRIMARY KEY,
  chain TEXT NOT NULL,
  name  TEXT NOT NULL,
  vm    TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_raw_message_topics_chain
  ON raw_message_topics(chain);

CREATE TABLE raw_messages (
  id           SERIAL PRIMARY KEY,
  topic_id     INTEGER NOT NULL,
  data         TEXT NOT NULL,
  hash         TEXT NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  processed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_raw_messages_topic
  ON raw_messages(topic_id);

CREATE UNIQUE INDEX idx_raw_messages_hash
  ON raw_messages(hash);

CREATE INDEX idx_raw_messages_created_at
  ON raw_messages(created_at);

CREATE INDEX idx_raw_messages_unprocessed
  ON raw_messages(processed_at) WHERE processed_at IS NULL;

-- +goose Down
DROP TABLE raw_message_topics;
DROP TABLE raw_messages;