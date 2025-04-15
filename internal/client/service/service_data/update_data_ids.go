package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) UpdateDataIDsOnServer(ctx context.Context, newDataIDs map[int]int) (err error) {
	action := "update data IDsOnServer"
	err = sv.rp.UpdateDataIDsOnServer(ctx, newDataIDs)
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
