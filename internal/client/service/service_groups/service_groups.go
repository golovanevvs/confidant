package service_groups

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportGroups interface {
}

type IRepositoryGroups interface {
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
	GetGroups(ctx context.Context, email string) (groups []model.Group, err error)
	GetGroupID(ctx context.Context, accountID int, titleGroup string) (groupID int, err error)
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
