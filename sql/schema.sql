CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    passhash text NOT NULL,
    salt text NOT NULL,
    phone_number text
);

CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    source text NOT NULL,
    destination text NOT NULL
    );

CREATE INDEX IF NOT EXISTS index_name
    ON links(source);