package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) GetGroupsByIDs(ctx context.Context, groupIDs []int) (groups []model.Group, err error) {
	action := "get groups by IDs"

	groups, err = sv.rp.GetGroupsByIDs(ctx, groupIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return groups, nil
}
