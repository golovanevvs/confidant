package service

import (
	"github.com/golovanevvs/confidant/internal/client/service/service_account"
	"github.com/golovanevvs/confidant/internal/client/service/service_groups"
	"github.com/golovanevvs/confidant/internal/client/service/service_manage"
	"github.com/golovanevvs/confidant/internal/client/service/service_security"
)

type ITransport interface {
	service_account.ITransportAccount
	service_manage.ITransportManage
}

type IRepository interface {
	service_account.IRepositoryAccount
	service_manage.IRepositoryManage
	service_groups.IRepositoryGroups
}

type service struct {
	*service_account.ServiceAccount
	*service_manage.ServiceManage
	*service_groups.ServiceGroups
	*service_security.ServiceSecurity
}

func New(tr ITransport, rp IRepository) *service {
	return &service{
		ServiceAccount:  service_account.New(tr, rp),
		ServiceManage:   service_manage.New(tr, rp),
		ServiceGroups:   service_groups.New(tr, rp),
		ServiceSecurity: service_security.New(),
	}
}
