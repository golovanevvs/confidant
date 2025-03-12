package service

import (
	"github.com/golovanevvs/confidant/internal/client/service/service_account"
	"github.com/golovanevvs/confidant/internal/client/service/service_manage"
)

type ITransport interface {
	service_account.ITransportAccount
	service_manage.ITransportManage
}

type IRepository interface {
	service_account.IRepositoryAccount
	service_manage.IRepositoryManage
}

type service struct {
	*service_account.ServiceAccount
	*service_manage.ServiceManage
}

func New(tr ITransport, rp IRepository) *service {
	return &service{
		ServiceAccount: service_account.New(tr, rp),
		ServiceManage:  service_manage.New(tr, rp),
	}
}
