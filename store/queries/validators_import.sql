INSERT INTO validators (
  node_id,
  stake_amount,
  stake_percent,
  potential_reward,
  reward_address,
  active,
  active_start_time,
  active_end_time,
  active_progress_percent,
  uptime,
  delegations_count,
  delegations_percent,
  delegated_amount,
  delegated_amount_percent,
  delegation_fee,
  capacity,
  capacity_percent,
  first_height,
  last_height,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (node_id) DO UPDATE
SET
  stake_amount             = excluded.stake_amount,
  stake_percent            = excluded.stake_percent,
  potential_reward         = excluded.potential_reward,
  reward_address           = excluded.reward_address,
  active                   = excluded.active,
  active_start_time        = excluded.active_start_time,
  active_end_time          = excluded.active_end_time,
  active_progress_percent  = excluded.active_progress_percent,
  uptime                   = excluded.uptime,
  delegations_count        = excluded.delegations_count,
  delegations_percent      = excluded.delegations_percent,
  delegated_amount         = excluded.delegated_amount,
  delegated_amount_percent = excluded.delegated_amount_percent,
  delegation_fee           = excluded.delegation_fee,
  capacity                 = excluded.capacity,
  capacity_percent         = excluded.capacity_percent,
  first_height             = excluded.first_height,
  last_height              = excluded.last_height,
  updated_at               = excluded.updated_at
