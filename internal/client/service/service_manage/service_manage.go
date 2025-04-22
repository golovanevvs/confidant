package service_manage

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportManage interface {
	GetServerStatus(ctx context.Context) (statusResp *model.TrResponse, err error)
}

type IRepositoryManage interface {
	CloseDB(ctx context.Context) error
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
