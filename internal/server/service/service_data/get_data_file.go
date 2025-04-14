package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataFile(ctx context.Context, dataID int) (file []byte, err error) {
	action := "get file by data ID"

	file, err = sv.rp.GetDataFile(ctx, dataID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.GroupsServiceErr,
			action,
			err,
		)
	}

	return file, nil
}
