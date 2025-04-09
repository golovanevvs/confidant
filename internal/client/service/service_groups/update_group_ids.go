package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) UpdateGroupIDsOnServer(ctx context.Context, newGroupIDs map[int]int) (err error) {
	action := "update group IDsOnServer"
	err = sv.rp.UpdateGroupIDsOnServer(ctx, newGroupIDs)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}
	return nil
}
