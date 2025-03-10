package service

import (
	"github.com/golovanevvs/confidant/internal/server/service/accountservice"
	"github.com/golovanevvs/confidant/internal/server/service/manageservice"
	"github.com/golovanevvs/confidant/internal/server/service/otherservice"
)

type IRepository interface {
	accountservice.IAccountRepository
	manageservice.IManageRepository
	otherservice.IManageRepository
}

type Service struct {
	*accountservice.AccountService
	*manageservice.ManageService
	*otherservice.OtherService
}

func New(rp IRepository) *Service {
	return &Service{
		AccountService: accountservice.New(rp),
		ManageService:  manageservice.New(rp),
		OtherService:   otherservice.New(rp),
	}
}
