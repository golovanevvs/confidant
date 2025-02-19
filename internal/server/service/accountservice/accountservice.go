package accountservice

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/repository"
)

type accountService struct {
	accountRp repository.IAccountRepository
}

func NewAccountService(accountRp repository.IAccountRepository) *accountService {
	return &accountService{
		accountRp: accountRp,
	}
}

func (sv *accountService) CreateUser(ctx context.Context) (int, error) {
	return 0, nil
}
