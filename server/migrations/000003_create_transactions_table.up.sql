CREATE TYPE payment AS (
    amount double precision,
    payer text
);

CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    payments payment[] NOT NULL,
    group_id bigint NOT NULL,
    version integer NOT NULL DEFAULT 1,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);