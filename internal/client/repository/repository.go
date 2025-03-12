package repository

import "github.com/golovanevvs/confidant/internal/client/repository/db_sqlite"

type Repository struct {
	*db_sqlite.SQLite
}

func New() (*Repository, error) {
	sqlite, err := db_sqlite.New()
	if err != nil {
		return nil, err
	}

	return &Repository{
		SQLite: sqlite,
	}, nil
}
