INSERT INTO validator_stats (
  node_id,
  time,
  bucket,
  uptime_min,
  uptime_max,
  uptime_avg,
  stake_amount,
  stake_percent,
  delegations_count,
  delegations_percent,
  delegated_amount,
  delegated_amount_percent
)
SELECT
  node_id,
  DATE_TRUNC('@bucket', time),
  '@bucket' AS bucket,
  ROUND(MIN(uptime), 4),
  ROUND(MAX(uptime), 4),
  ROUND(AVG(uptime), 4),
  AVG(stake_amount),
  AVG(stake_percent),
  ROUND(AVG(delegations_count)),
  ROUND(AVG(delegations_percent), 2),
  AVG(delegated_amount),
  ROUND(AVG(delegated_amount_percent), 2)
FROM
  validator_sequences
WHERE
  time >= ? AND time <= ?
GROUP BY
  node_id,
  DATE_TRUNC('@bucket', time)
