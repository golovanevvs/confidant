package repository

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/server/repository/postgres"
)

type IUserRepository interface {
	SaveUser(ctx context.Context) (int, error)
}

type IMyRepository interface {
	DoIt(ctx context.Context) (int, error)
}

type IYourRepository interface {
	CrashIt(ctx context.Context) (int, error)
}

type Repository struct {
	IUserRepository
	IMyRepository
	IYourRepository
}

func New(databaseURI string) (*Repository, error) {
	db, err := postgres.New(databaseURI)
	if err != nil {
		return nil, fmt.Errorf("postgres DB initialization error: %s", err.Error())
	}
	return &Repository{
		IUserRepository: postgres.NewUserPostgres(db),
		IMyRepository:   postgres.NewMyPostgres(db),
		IYourRepository: postgres.NewYourPostgres(db),
	}, nil
}
