package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) GetEmails(ctx context.Context, groupIDs []int) (mapGroupIDEmails map[int][]string, err error) {
	action := "get map groupID-emails by group IDs"

	mapGroupIDEmails, err = sv.rp.GetEmails(ctx, groupIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.GroupsServiceErr,
			action,
			err,
		)
	}

	return mapGroupIDEmails, nil
}
