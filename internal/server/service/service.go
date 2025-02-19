package service

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/service/accountservice"
)

type IAuthService interface {
	CreateUser(ctx context.Context) (int, error)
}

type IMyService interface {
	MyTask(ctx context.Context) (int, error)
}

type IYourService interface {
	YourTask(ctx context.Context) (int, error)
}

type Service struct {
	IAuthService
	IMyService
	IYourService
}

func New(rp *repository.Repository) *Service {
	return &Service{
		IAuthService: accountservice.NewAccountService(rp.IAccountRepository),
		IMyService:   newMyService(rp.IMyRepository),
		IYourService: newYourService(rp.IYourRepository),
	}
}
