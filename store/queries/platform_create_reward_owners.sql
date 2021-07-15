INSERT INTO rewards_owners (id, locktime, threshold)
VALUES (?, ?, ?)
ON CONFLICT (id) DO NOTHING