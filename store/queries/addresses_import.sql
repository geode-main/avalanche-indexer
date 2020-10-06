INSERT INTO addresses (
  value,
  balance,
  unlocked_balance,
  locked_not_stakeable,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (value) DO UPDATE
SET
  balance              = excluded.balance,
  unlocked_balance     = excluded.unlocked_balance,
  locked_not_stakeable = excluded.locked_not_stakeable,
  updated_at           = excluded.updated_at
