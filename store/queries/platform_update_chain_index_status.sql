INSERT INTO sync_statuses (id, index_id, index_time, tip_id, tip_time)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE
SET
  index_id   = excluded.index_id,
  index_time = excluded.index_time,
  tip_id     = excluded.tip_id,
  tip_time   = excluded.tip_time