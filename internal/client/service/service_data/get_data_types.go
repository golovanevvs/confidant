package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataTypes(ctx context.Context, accountID int, groupID int) (dataTypes []string, err error) {
	action := "get data types"

	dataTypes, err = sv.rp.GetDataTypes(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return dataTypes, nil
}
