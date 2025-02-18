package service

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/repository"
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

type service struct {
	IAuthService
	IMyService
	IYourService
}

func New(rp *repository.Repository) *service {
	return &service{
		IAuthService: newAuthService(rp.IUserRepository),
		IMyService:   newMyService(rp.IMyRepository),
		IYourService: newYourService(rp.IYourRepository),
	}
}
