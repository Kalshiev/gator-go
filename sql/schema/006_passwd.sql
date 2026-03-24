-- +goose up
ALTER TABLE users ADD COLUMN password_hash TEXT NOT NULL;

-- +goose down
ALTER TABLE users DROP COLUMN password_hash;