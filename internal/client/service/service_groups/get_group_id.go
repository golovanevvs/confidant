package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) GetGroupID(ctx context.Context, accountID int, titleGroup string) (groupID int, err error) {
	action := "get group ID"

	groupID, err = sv.rp.GetGroupID(ctx, accountID, titleGroup)
	if err != nil {
		return -1, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return groupID, nil
}
