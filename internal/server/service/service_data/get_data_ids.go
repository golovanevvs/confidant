package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error) {
	action := "get data IDs"

	dataIDs, err = sv.rp.GetDataIDs(ctx, accountID)
	if err != nil {
		return dataIDs, fmt.Errorf(
			"%s: %s: %w",
			customerrors.DataServiceErr,
			action,
			err,
		)
	}

	return dataIDs, nil
}
