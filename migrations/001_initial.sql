CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE TABLE IF NOT EXISTS dictionary  (
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL,
    definition TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_word_trgm ON dictionary USING gin(word gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_word_fts ON dictionary USING gin( to_tsvector('russian', word));