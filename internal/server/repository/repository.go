package repository

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/golovanevvs/confidant/internal/server/repository/postgres"
)

type IManageRepository interface {
	CloseDB() error
}

type IAccountRepository interface {
	SaveAccount(ctx context.Context, account model.Account) (int, error)
	LoadAccountID(ctx context.Context, email, passwordHash string) (int, error)
}

type IMyRepository interface {
	DoIt(ctx context.Context) (int, error)
}

type IYourRepository interface {
	CrashIt(ctx context.Context) (int, error)
}

type Repository struct {
	IManageRepository
	IAccountRepository
	IMyRepository
	IYourRepository
}

func New(databaseURI string) (*Repository, error) {
	db, err := postgres.New(databaseURI)
	if err != nil {
		return nil, fmt.Errorf("postgres DB initialization error: %s", err.Error())
	}
	return &Repository{
		IManageRepository:  postgres.NewManagePostgres(db),
		IAccountRepository: postgres.NewAccountPostgres(db),
		IMyRepository:      postgres.NewMyPostgres(db),
		IYourRepository:    postgres.NewYourPostgres(db),
	}, nil
}
