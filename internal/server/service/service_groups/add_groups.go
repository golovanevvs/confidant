package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (sv *ServiceGroups) AddGroups(ctx context.Context, groups []model.Group) (groupIDs map[int]int, err error) {
	action := "add groups"

	groupIDs, err = sv.rp.AddGroups(ctx, groups)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.GroupsServiceErr,
			action,
			err,
		)
	}

	return groupIDs, nil
}
