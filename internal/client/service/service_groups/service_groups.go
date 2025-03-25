package service_groups

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportGroups interface {
}

type IRepositoryGroups interface {
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
}

type ServiceGroups struct {
	tr ITransportGroups
	rp IRepositoryGroups
}

func New(tr ITransportGroups, rp IRepositoryGroups) *ServiceGroups {
	return &ServiceGroups{
		tr: tr,
		rp: rp,
	}
}
