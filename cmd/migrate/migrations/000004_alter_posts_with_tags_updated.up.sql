ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS tags varchar(100)[] NOT NULL DEFAULT '{}';

ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS updated_at timestamptz(0) NOT NULL DEFAULT now();