package service_account

import (
	"context"
)

type IRepositoryGroups interface {
	GetGroupIDs(ctx context.Context, accountID int) (groupIDs map[int]struct{}, err error)
}

type ServiceGroups struct {
	rp IRepositoryGroups
}

func New(groupsRp IRepositoryGroups) *ServiceGroups {
	return &ServiceGroups{
		rp: groupsRp,
	}
}
