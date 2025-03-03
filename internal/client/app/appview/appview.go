package appview

import (
	"go.uber.org/zap"
)

type IAccountService interface {
	Register(email, password string) (accountID int, err error)
	// Login(email, password string) (accountID int, err error)
	// ChangePassword(email, password, newPassword string) error
	//GetUser(login string) (string, error)
}

type IService interface {
	IAccountService
}

type AppView struct {
	sv IService
	lg *zap.SugaredLogger
}

func New(sv IService, lg *zap.SugaredLogger) *AppView {
	return &AppView{
		sv: sv,
		lg: lg,
	}
}
