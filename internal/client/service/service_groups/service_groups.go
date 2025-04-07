package service_groups

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportGroups interface {
	GetGroupIDs(ctx context.Context, accessToken string) (groupIDs map[int]struct{}, err error)
}

type IRepositoryGroups interface {
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
	GetGroups(ctx context.Context, email string) (groups []model.Group, err error)
	GetGroupID(ctx context.Context, email string, titleGroup string) (groupID int, err error)
	AddEmail(ctx context.Context, groupID int, email string) (err error)
	GetGroupIDs(ctx context.Context, email string) (groupServerIDs map[int]struct{}, groupNoServerIDs map[int]struct{}, err error)
	AddGroupBySync(ctx context.Context, group model.Group) (err error)
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
