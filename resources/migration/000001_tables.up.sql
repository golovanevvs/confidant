BEGIN;

CREATE TABLE IF NOT EXISTS account(
    id SERIAL PRIMARY KEY,
    email VARCHAR(250) NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_token(
    id SERIAL PRIMARY KEY,
    account_id INT,
    token_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS groups(
    id SERIAL PRIMARY KEY,
	id_on_client INT,
    title VARCHAR(250),
    account_id INT
);

CREATE TABLE IF NOT EXISTS email_in_groups(
    id SERIAL PRIMARY KEY,
    group_id INT,
    email VARCHAR(250) NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS data(
	id SERIAL PRIMARY KEY,
	id_on_client INT,
	group_id INT,
	data_type TEXT NOT NULL,
	title TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS data_note(
    id SERIAL PRIMARY KEY,
	data_id INT,
	descr BYTEA,
	note BYTEA NOT NULL,
	FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS data_pass(
	id SERIAL PRIMARY KEY,
	data_id INT,
	descr BYTEA,
	login BYTEA NOT NULL,
	pass BYTEA,
	FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS data_card(
	id SERIAL PRIMARY KEY,
	data_id INT,
	descr BYTEA,
	number BYTEA NOT NULL,
	date BYTEA,
	name BYTEA,
	cvc2 BYTEA,
	pin BYTEA,
	bank BYTEA,
	FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS data_file(
	id SERIAL PRIMARY KEY,
	data_id INT,
	descr BYTEA,
	filename BYTEA NOT NULL,
	filesize BYTEA,
	filedate BYTEA,
	file BYTEA,
	FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
);

COMMIT;