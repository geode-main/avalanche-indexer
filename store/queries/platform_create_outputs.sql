INSERT INTO transaction_outputs (
  id,
  tx_id,
  chain,
  asset,
  type,
  index,
  locktime,
  threshold,
  amount,
  group_id,
  stake,
  reward,
  spent,
  spent_tx_id,
  addresses,
  payload
)
VALUES @values

ON CONFLICT (id) DO NOTHING