package repository

import "github.com/golovanevvs/confidant/internal/client/repository/sqlite"

type Repository struct {
	*sqlite.SQLite
}

func New() (*Repository, error) {
	sqlite, err := sqlite.New()
	if err != nil {
		return nil, err
	}

	return &Repository{
		SQLite: sqlite,
	}, nil
}
