package service_data

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataDates(ctx context.Context, dataIDs []int) (dataDatesFromClient map[int]time.Time, err error) {
	action := "get dates by data IDs"

	dataDatesFromClient, err = sv.rp.GetDataDates(ctx, dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return dataDatesFromClient, nil
}
