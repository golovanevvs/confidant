package service_groups

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportGroups interface {
	GetGroupIDs(ctx context.Context, accessToken string) (trResponse *model.GroupSyncResp, err error)
}

type IRepositoryGroups interface {
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
	GetGroups(ctx context.Context, email string) (groups []model.Group, err error)
	GetGroupsByIDs(ctx context.Context, groupIDs []int) (groups []model.Group, err error)
	GetGroupID(ctx context.Context, email string, titleGroup string) (groupID int, err error)
	AddEmail(ctx context.Context, groupID int, email string) (err error)
	GetGroupIDs(ctx context.Context, email string) (groupServerIDs []int, groupNoServerIDs []int, err error)
	AddGroupBySync(ctx context.Context, group model.Group) (err error)
	UpdateGroupIDsOnServer(ctx context.Context, newGroupIDs map[int]int) (err error)
	GetEmails(ctx context.Context, groupIDs []int) (mapGroupIDEmails map[int][]string, err error)
	AddEmailsBySync(ctx context.Context, mapGroupIDEmails map[int][]string) (err error)
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
