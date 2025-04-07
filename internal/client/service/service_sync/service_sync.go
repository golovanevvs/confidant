package service_sync

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type IServiceGroups interface {
	GetGroups(ctx context.Context, email string) (groups []model.Group, err error)
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
	GetGroupID(ctx context.Context, email string, titleGroup string) (groupID int, err error)
	AddEmail(ctx context.Context, groupID int, email string) (err error)
	GetGroupIDs(ctx context.Context, email string) (groupServerIDs map[int]struct{}, groupNoServerIDs map[int]struct{}, err error)
	AddGroupBySync(ctx context.Context, group model.Group) (err error)
}

type ITransportSync interface {
	GetGroupIDs(ctx context.Context, accessToken string) (groupIDs map[int]struct{}, err error)
	GetGroups(ctx context.Context, accessToken string, groupIDs map[int]struct{}) (groupsFromServer []model.Group, err error)
}

type ServiceSync struct {
	tr ITransportSync
	sg IServiceGroups
}

func New(tr ITransportSync, sg IServiceGroups) *ServiceSync {
	return &ServiceSync{
		tr: tr,
		sg: sg,
	}
}
