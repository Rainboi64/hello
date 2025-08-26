CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    passhash text NOT NULL,
    salt text NOT NULL,
    phone_number text
);

CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    destination text NOT NULL,

    createdBy BIGSERIAL NOT NULL,
    FOREIGN KEY (createdBy) REFERENCES users(id)
)