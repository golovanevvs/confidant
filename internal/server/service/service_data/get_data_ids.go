package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error) {
	//! ----------------- СТОП --------------------
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
