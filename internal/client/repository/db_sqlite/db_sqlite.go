package db_sqlite

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	*sqliteManage
	*sqliteAccount
	*sqliteGroups
}

func New() (*SQLite, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "confidant_client.db")

	var db *sqlx.DB
	ctx := context.Background()

	_, err = os.Stat(dbFile)
	if err != nil {
		db, err = sqlx.Open("sqlite3", dbFile)
		if err != nil {
			return nil, err
		}
		_, err = db.ExecContext(ctx, `

			CREATE TABLE IF NOT EXISTS account(
    			id INTEGER PRIMARY KEY,
    			email TEXT NOT NULL UNIQUE,
    			password_hash BLOB NOT NULL
			);

			CREATE TABLE IF NOT EXISTS active_account(
				account_id INTEGER,
				token TEXT,
				dummy INTEGER DEFAULT 1 UNIQUE
			);

			CREATE TABLE IF NOT EXISTS groups(
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
				id_on_server INTEGER,
   				title TEXT NOT NULL,
    			account_id INTEGER,
    			FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS emails_in_groups(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				email TEXT NOT NULL,
				groups_id INTEGER,
				FOREIGN KEY (groups_id) REFERENCES groups (id) ON DELETE CASCADE
			);

		`)
		if err != nil {
			return nil, err
		}

		_, err = db.ExecContext(ctx, `

			CREATE INDEX email_index
			ON account (email);

		`)
		if err != nil {
			return nil, err
		}
	} else {
		db, err = sqlx.Open("sqlite3", dbFile)
		if err != nil {
			return nil, err
		}
	}

	return &SQLite{
		sqliteManage:  NewSQLiteManage(db),
		sqliteAccount: NewSQLiteAccount(db),
		sqliteGroups:  NewSQLiteGroups(db),
	}, nil
}
