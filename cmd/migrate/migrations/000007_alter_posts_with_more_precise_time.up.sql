ALTER TABLE posts
    ALTER COLUMN created_at TYPE timestamp(6) with time zone;

ALTER TABLE posts
    ALTER COLUMN updated_at TYPE timestamp(6) with time zone;
