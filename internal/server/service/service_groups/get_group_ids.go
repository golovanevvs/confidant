package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) GetGroupIDs(ctx context.Context, accountID int) (groupIDs map[int]struct{}, err error) {
	action := "get group IDs"

	groupIDs, err = sv.rp.GetGroupIDs(ctx, accountID)
	if err != nil {
		return groupIDs, fmt.Errorf(
			"%s: %s: %w",
			customerrors.GroupsServiceErr,
			action,
			err,
		)
	}

	return groupIDs, nil
}
