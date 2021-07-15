INSERT INTO raw_messages (
  topic_id,
  data,
  hash,
  created_at
)
VALUES (?, ?, ?, ?)
ON CONFLICT (hash) DO NOTHING