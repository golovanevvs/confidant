BEGIN;

CREATE TABLE account(
    id SERIAL PRIMARY KEY,
    email VARCHAR(250) NOT NULL UNIQUE,
    password_hash VARCHAR(250) NOT NULL
);

CREATE TABLE groups(
    id SERIAL PRIMARY KEY,
    title VARCHAR(250),
    account_id INT,
    FOREIGN KEY (account_id) REFERENCES account (id)
);

COMMIT;