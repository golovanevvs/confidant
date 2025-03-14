BEGIN;

CREATE TABLE IF NOT EXISTS account(
    id SERIAL PRIMARY KEY,
    email VARCHAR(250) NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens(
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES account(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS groups(
    id SERIAL PRIMARY KEY,
    title VARCHAR(250),
    account_id INT,
    FOREIGN KEY (account_id) REFERENCES account (id)
);

CREATE TABLE IF NOT EXISTS emails_of_groups(
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    email VARCHAR(250) NOT NULL
);

COMMIT;