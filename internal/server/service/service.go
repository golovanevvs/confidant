package service

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/service/accountservice"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, account model.Account) (int, error)
	BuildJWTString(ctx context.Context, accountID int) (string, error)
}

// type IMyService interface {
// 	MyTask(ctx context.Context) (int, error)
// }

// type IYourService interface {
// 	YourTask(ctx context.Context) (int, error)
// }

type Service struct {
	IAccountService
	//IMyService
	//IYourService
}

func New(rp *repository.Repository) *Service {
	return &Service{
		IAccountService: accountservice.NewAccountService(rp.IAccountRepository),
		// IMyService:      newMyService(rp.),
		// IYourService:    newYourService(rp.IYourRepository),
	}
}
