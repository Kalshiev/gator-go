-- +goose up
ALTER TABLE feeds
ADD COLUMN IF NOT EXISTS last_fetched_at TIMESTAMP NULL;

-- +goose down
ALTER TABLE feeds
DROP COLUMN IF EXISTS last_fetched_at;