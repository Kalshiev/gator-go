-- +goose up
CREATE TABLE IF NOT EXISTS feeds_follow (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (feed_id, user_id)
);

-- +goose down
DROP TABLE IF EXISTS feeds_follow;