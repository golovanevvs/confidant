package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (sv *ServiceGroups) GetGroups(ctx context.Context, accountID int, groupIDs []int) (groups []model.Group, err error) {
	action := "get groups"

	groups, err = sv.rp.GetGroups(ctx, accountID, groupIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.GroupsServiceErr,
			action,
			err,
		)
	}

	return groups, nil
}
