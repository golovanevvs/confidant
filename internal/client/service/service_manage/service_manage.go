package service_manage

import (
	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportManage interface {
	ServerStatus() (statusResp *model.TrResponse, err error)
}

type IRepositoryManage interface {
	GetServerStatus() (statusResp *model.StatusResp, err error)
}

type ServiceManage struct {
	tr ITransportManage
	rp IRepositoryManage
}

func New(tr ITransportManage, rp IRepositoryManage) *ServiceManage {
	return &ServiceManage{
		tr: tr,
		rp: rp,
	}
}
