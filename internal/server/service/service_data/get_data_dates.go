package service_data

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataDates(ctx context.Context, dataIDs []int) (datadates map[int]time.Time, err error) {
	action := "get dates by data IDs"

	datadates, err = sv.rp.GetDataDates(ctx, dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.GroupsServiceErr,
			action,
			err,
		)
	}

	return datadates, nil
}
