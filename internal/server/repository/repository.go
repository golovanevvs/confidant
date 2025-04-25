package repository

import (
	"github.com/golovanevvs/confidant/internal/server/repository/postgres"
)

type Repository struct {
	*postgres.Postgres
}

func New(databaseURI string) (*Repository, error) {
	postgres, err := postgres.New(databaseURI)
	if err != nil {
		return nil, err
	}

	return &Repository{
		postgres,
	}, nil
}
