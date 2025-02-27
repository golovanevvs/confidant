package repository

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IManageRepository interface {
	CloseDB() error
}

type IAccountRepository interface {
	SaveAccount(ctx context.Context, account model.Account) (int, error)
	LoadAccountID(ctx context.Context, email, passwordHash string) (int, error)
}

// type IMyRepository interface {
// 	DoIt(ctx context.Context) (int, error)
// }

// type IYourRepository interface {
// 	CrashIt(ctx context.Context) (int, error)
// }

type Repository struct {
	IManageRepository
	IAccountRepository
	// IMyRepository
	// IYourRepository
}

func New(
	IManageRepository IManageRepository,
	IAccountRepository IAccountRepository,
	// IMyRepository IMyRepository,
	// IYourRepository IYourRepository,
) *Repository {

	return &Repository{
		IManageRepository:  IManageRepository,
		IAccountRepository: IAccountRepository,
		// IMyRepository:      IMyRepository,
		// IYourRepository:    IYourRepository,
	}
}
