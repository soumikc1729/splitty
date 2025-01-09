CREATE TABLE IF NOT EXISTS groups (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    token text NOT NULL UNIQUE,
    users text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);