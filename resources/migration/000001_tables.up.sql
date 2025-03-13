BEGIN;

CREATE TABLE account(
    id SERIAL PRIMARY KEY,
    email VARCHAR(250) NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES account(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE,
);

CREATE TABLE groups(
    id SERIAL PRIMARY KEY,
    title VARCHAR(250),
    account_id INT,
    FOREIGN KEY (account_id) REFERENCES account (id)
);

COMMIT;