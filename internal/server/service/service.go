package service

import (
	"github.com/golovanevvs/confidant/internal/server/service/accountservice"
	"github.com/golovanevvs/confidant/internal/server/service/otherservice"
)

type IRepository interface {
	accountservice.IAccountRepository
	otherservice.IManageRepository
}

type Service struct {
	*accountservice.AccountService
	*otherservice.OtherService
}

func New(rp IRepository) *Service {
	return &Service{
		AccountService: accountservice.New(rp),
		OtherService:   otherservice.New(rp),
	}
}
