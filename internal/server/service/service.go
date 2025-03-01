package service

import (
	"github.com/golovanevvs/confidant/internal/server/repository"
	"github.com/golovanevvs/confidant/internal/server/service/accountservice"
	"github.com/golovanevvs/confidant/internal/server/transport/http/handler"
)

// type IAccountService interface {
// 	CreateAccount(ctx context.Context, account model.Account) (int, error)
// 	BuildJWTString(ctx context.Context, accountID int) (string, error)
// }

// type IMyService interface {
// 	MyTask(ctx context.Context) (int, error)
// }

// type IYourService interface {
// 	YourTask(ctx context.Context) (int, error)
// }

// type Service struct {
// 	handler.IAccountService
// 	// IAccountService
// 	//IMyService
// 	//IYourService
// }

func New(rp *repository.Repository) *handler.Service {
	return &handler.Service{
		IAccountService: accountservice.NewAccountService(rp.IAccountRepository),
		// IMyService:      newMyService(rp.),
		// IYourService:    newYourService(rp.IYourRepository),
	}
}
