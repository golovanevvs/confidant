package sqlite

import (
	"context"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type SQLite struct {
	*sqliteAccount
	*sqliteManage
}

func New() (*SQLite, error) {
	// db, err := sqlx.Open("sqlite3", databasePath)
	// if err != nil {
	// 	return nil, err
	// }
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

		`)
		if err != nil {
			return nil, err
		}

		_, err = db.ExecContext(ctx, `

			CREATE INDEX email
			ON confidant_client (email);

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
		sqliteAccount: NewSQLiteAccount(db),
		sqliteManage:  NewSQLiteManage(db),
	}, nil
}
