-- +goose Up
ALTER TABLE network_metrics ADD COLUMN total_staked DECIMAL;
ALTER TABLE network_metrics ADD COLUMN total_delegated DECIMAL;

ALTER TABLE network_stats ADD COLUMN total_staked DECIMAL;
ALTER TABLE network_stats ADD COLUMN total_delegated DECIMAL;

-- +goose Down
ALTER TABLE network_metrics DROP COLUMN total_staked;
ALTER TABLE network_metrics DROP COLUMN total_delegated;

ALTER TABLE network_stats DROP COLUMN total_staked;
ALTER TABLE network_stats DROP COLUMN total_delegated;