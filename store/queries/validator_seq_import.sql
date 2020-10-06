INSERT INTO validator_sequences (
  time,
  height,
  node_id,
  stake_amount,
  stake_percent,
  potential_reward,
  reward_address,
  active,
  active_start_time,
  active_end_time,
  active_progress_percent,
  delegations_count,
  delegations_percent,
  delegated_amount,
  delegated_amount_percent,
  delegation_fee,
  uptime
)
VALUES @values
