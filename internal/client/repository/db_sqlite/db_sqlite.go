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
	*sqliteData
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
		
			PRAGMA foreign_keys = ON
		
		`)
		if err != nil {
			return nil, err
		}

		_, err = db.ExecContext(ctx, `

			CREATE TABLE IF NOT EXISTS account(
    			id INTEGER PRIMARY KEY,
    			email TEXT NOT NULL UNIQUE,
    			password_hash BLOB NOT NULL,
				token TEXT
			);

			CREATE TABLE IF NOT EXISTS active_account(
				account_id INTEGER,
				dummy INTEGER DEFAULT 1 UNIQUE
			);

			CREATE TABLE IF NOT EXISTS groups(
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
				id_on_server INTEGER DEFAULT -1,
   				title TEXT NOT NULL,
    			account_id INTEGER,
    			FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS email_in_groups(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				email TEXT NOT NULL,
				group_id INTEGER,
				FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS data(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				id_on_server INTEGER DEFAULT -1,
				group_id INTEGER,
				data_type TEXT NOT NULL,
				title TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
			);
			
			CREATE TABLE IF NOT EXISTS data_note(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				data_id INTEGER,
				desc BLOB,
				note BLOB NOT NULL,
				FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS data_pass(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				data_id INTEGER,
				desc BLOB,
				login BLOB NOT NULL,
				pass BLOB,
				FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS data_card(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				data_id INTEGER,
				desc BLOB,
				number BLOB NOT NULL,
				date BLOB,
				name BLOB,
				cvc2 BLOB,
				pin BLOB,
				bank BLOB,
				FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS data_file(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				data_id INTEGER,
				desc BLOB,
				filename BLOB NOT NULL,
				filesize BLOB,
				filedate BLOB,
				file BLOB,
				FOREIGN KEY (data_id) REFERENCES data (id) ON DELETE CASCADE
			);

		`)
		if err != nil {
			return nil, err
		}

		_, err = db.ExecContext(ctx, `
		
			CREATE INDEX group_id_index
			ON data (group_id);
		
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
		sqliteData:    NewSQLiteData(db),
	}, nil
}
