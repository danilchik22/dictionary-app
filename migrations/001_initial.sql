CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE TABLE IF NOT EXISTS dictionary  (
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL,
    definition TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email_lower_pattern ON dictionary (lower(word) varchar_pattern_ops);