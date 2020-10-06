INSERT INTO delegations (
  reference_id,
  node_id,
  stake_amount,
  potential_reward,
  reward_address,
  active,
  active_start_time,
  active_end_time,
  first_height,
  last_height,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (reference_id) DO UPDATE
SET
  node_id           = excluded.node_id,
  stake_amount      = excluded.stake_amount,
  potential_reward  = excluded.potential_reward,
  reward_address    = excluded.reward_address,
  active            = excluded.active,
  active_start_time = excluded.active_start_time,
  active_end_time   = excluded.active_end_time,
  last_height       = excluded.last_height,
  updated_at        = excluded.updated_at
