INSERT INTO rewards_owner_outputs (id, transaction_id, index)
VALUES (?, ?, ?)
ON CONFLICT (id) DO NOTHING