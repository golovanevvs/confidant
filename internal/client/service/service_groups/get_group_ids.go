package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) GetGroupIDs(ctx context.Context, email string) (groupServerIDs []int, groupNoServerIDs []int, err error) {
	action := "get group IDs"

	groupServerIDs, groupNoServerIDs, err = sv.rp.GetGroupIDs(ctx, email)
	if err != nil {
		return groupServerIDs, groupNoServerIDs, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return groupServerIDs, groupNoServerIDs, nil
}
