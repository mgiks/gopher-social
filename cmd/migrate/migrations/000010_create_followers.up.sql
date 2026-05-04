CREATE TABLE IF NOT EXISTS followers (
    user_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    follower_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone DEFAULT now(),
    PRIMARY KEY (user_id, follower_id)
);
