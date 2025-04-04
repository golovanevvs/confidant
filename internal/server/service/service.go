package service

import (
	"github.com/golovanevvs/confidant/internal/server/service/service_account"
	"github.com/golovanevvs/confidant/internal/server/service/service_groups"
	"github.com/golovanevvs/confidant/internal/server/service/service_manage"
)

type IRepository interface {
	service_account.IRepositoryAccount
	service_manage.IRepositoryManage
	service_groups.IRepositoryGroups
}

type service struct {
	*service_account.ServiceAccount
	*service_manage.ServiceManage
	*service_groups.ServiceGroups
}

func New(rp IRepository) *service {
	return &service{
		ServiceAccount: service_account.New(rp),
		ServiceManage:  service_manage.New(rp),
		ServiceGroups:  service_groups.New(rp),
	}
}
