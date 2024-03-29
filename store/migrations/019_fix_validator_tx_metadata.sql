-- +goose Up
UPDATE transactions
SET
  metadata = json_build_object(
    'node_id',         metadata->>'validator_node_id',
    'weight',          (metadata->>'validator_weight')::bigint,
    'start_time',      metadata->>'validator_start_time',
    'end_time',        metadata->>'validator_end_time',
    'duration',        (metadata->>'validator_duration')::bigint,
    'commission_rate', (metadata->>'validator_commission_rate')::bigint
  )
WHERE
  type = 'p_add_validator'
  AND metadata->>'validator_node_id' IS NOT NULL;

UPDATE transactions
SET
  metadata = json_build_object(
    'node_id',    metadata->>'validator_node_id',
    'weight',     (metadata->>'validator_weight')::bigint,
    'start_time', metadata->>'validator_start_time',
    'end_time',   metadata->>'validator_end_time',
    'duration',   (metadata->>'validator_duration')::bigint
  )
WHERE
  type = 'p_add_delegator'
  AND metadata->>'validator_node_id' IS NOT NULL;