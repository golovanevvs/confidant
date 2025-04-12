package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataIDs(ctx context.Context, groupIDs []int) (dataServerIDs []int, dataNoServerIDs []int, err error) {
	action := "get data IDs"

	dataServerIDs, dataNoServerIDs, err = sv.rp.GetDataIDs(ctx, groupIDs)
	if err != nil {
		return dataServerIDs, dataNoServerIDs, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return dataServerIDs, dataNoServerIDs, nil
}
