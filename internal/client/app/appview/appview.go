package appview

import (
	"github.com/golovanevvs/confidant/internal/client/model"
	"go.uber.org/zap"
)

type IAccountService interface {
	RegisterAccount(email, password string) (registerAccountResp *model.RegisterAccountResp, err error)
	// Login(email, password string) (accountID int, err error)
	// ChangePassword(email, password, newPassword string) error
	//GetUser(login string) (string, error)
}

type IStatusServerService interface {
	GetServerStatus() (statusResp *model.StatusResp, err error)
}

type IService interface {
	IAccountService
	IStatusServerService
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
