package service

import (
	"github.com/golovanevvs/confidant/internal/server/service/service_account"
	"github.com/golovanevvs/confidant/internal/server/service/service_manage"
)

type IRepository interface {
	service_account.IRepositoryAccount
	service_manage.IRepositoryManage
}

type service struct {
	*service_account.ServiceAccount
	*service_manage.ServiceManage
}

func New(rp IRepository) *service {
	return &service{
		ServiceAccount: service_account.New(rp),
		ServiceManage:  service_manage.New(rp),
	}
}
