package service

import (
	"github.com/golovanevvs/confidant/internal/client/service/service_account"
	"github.com/golovanevvs/confidant/internal/client/service/service_data"
	"github.com/golovanevvs/confidant/internal/client/service/service_groups"
	"github.com/golovanevvs/confidant/internal/client/service/service_manage"
	"github.com/golovanevvs/confidant/internal/client/service/service_sync"
)

type ITransport interface {
	service_account.ITransportAccount
	service_manage.ITransportManage
	service_groups.ITransportGroups
	service_sync.ITransportSync
}

type IRepository interface {
	service_account.IRepositoryAccount
	service_manage.IRepositoryManage
	service_groups.IRepositoryGroups
	service_data.IRepositoryData
}

type service struct {
	*service_account.ServiceAccount
	*service_manage.ServiceManage
	*service_groups.ServiceGroups
	*service_data.ServiceData
	*service_sync.ServiceSync
}

func New(tr ITransport, rp IRepository) *service {
	sg := service_groups.New(tr, rp)
	return &service{
		ServiceAccount: service_account.New(tr, rp),
		ServiceManage:  service_manage.New(tr, rp),
		ServiceGroups:  sg,
		ServiceData:    service_data.New(tr, rp),
		ServiceSync:    service_sync.New(tr, sg),
	}
}
