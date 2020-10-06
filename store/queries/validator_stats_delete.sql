DELETE FROM validator_stats
WHERE time::timestamp = ? AND bucket = ?
