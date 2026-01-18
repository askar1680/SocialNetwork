DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_user'
          AND conrelid = 'posts'::regclass
    ) THEN
ALTER TABLE posts
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);
END IF;
END $$;
