package service

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/repository"
)

type authService struct {
	userRp repository.IUserRepository
}

func newAuthService(userRp repository.IUserRepository) *authService {
	return &authService{
		userRp: userRp,
	}
}

func (sv *authService) CreateUser(ctx context.Context) (int, error) {
	return 0, nil
}
