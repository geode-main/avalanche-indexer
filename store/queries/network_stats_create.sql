INSERT INTO network_stats (
  time,
  bucket,
  height_change,
  peers,
  blockchains,
  active_validators,
  pending_validators,
  active_delegations,
  pending_delegations,
  uptime
)
SELECT
  DATE_TRUNC('@bucket', time),
  '@bucket',
  MAX(height) - MIN(height),
  AVG(peers_count),
  MAX(blockchains_count),
  AVG(active_validators_count),
  AVG(pending_validators_count),
  AVG(active_delegations_count),
  AVG(pending_delegations_count),
  ROUND(AVG(uptime), 4)
FROM
  network_metrics
WHERE
  time >= ? AND time <= ?
GROUP BY
  DATE_TRUNC('@bucket', time)

ON CONFLICT (time, bucket) DO UPDATE
SET
  height_change       = excluded.height_change,
  peers               = excluded.peers,
  blockchains         = excluded.blockchains,
  active_validators   = excluded.active_validators,
  pending_validators  = excluded.pending_validators,
  active_delegations  = excluded.active_delegations,
  pending_delegations = excluded.pending_delegations,
  uptime              = excluded.uptime
