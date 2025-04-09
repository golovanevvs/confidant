package service_groups

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IRepositoryGroups interface {
	GetGroupIDs(ctx context.Context, accountID int) (groupIDs map[int]struct{}, err error)
	GetGroups(ctx context.Context, accountID int, groupIDs []int) (groups []model.Group, err error)
	AddGroups(ctx context.Context, groups []model.Group) (groupIDs map[int]int, err error)
}

type ServiceGroups struct {
	rp IRepositoryGroups
}

func New(groupsRp IRepositoryGroups) *ServiceGroups {
	return &ServiceGroups{
		rp: groupsRp,
	}
}
