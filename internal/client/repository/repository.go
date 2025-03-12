package repository

import "github.com/golovanevvs/confidant/internal/client/repository/sqlite"

type Repository struct {
	*sqlite.SQLite
}

func New(databasePath string) (*Repository, error) {
	sqlite, err := sqlite.New(databasePath)
	if err != nil {
		return nil, err
	}

	return &Repository{
		SQLite: sqlite,
	}, nil
}
