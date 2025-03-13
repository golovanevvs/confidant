package db_sqlite

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	*sqliteManage
	*sqliteAccount
}

func New() (*SQLite, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "confidant_client.db")
	fmt.Println(dbFile)

	var db *sqlx.DB
	ctx := context.Background()

	_, err = os.Stat(dbFile)
	if err != nil {
		db, err = sqlx.Open("sqlite3", dbFile)
		if err != nil {
			return nil, err
		}
		_, err = db.ExecContext(ctx, `

			CREATE TABLE account(
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			email TEXT NOT NULL UNIQUE,
    			password_hash BLOB NOT NULL,
				refresh_token TEXT
			);



			CREATE TABLE groups(
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
   				title TEXT,
    			account_id INTEGER,
    			FOREIGN KEY (account_id) REFERENCES account (id)
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
	}, nil
}
