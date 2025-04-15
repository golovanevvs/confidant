package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataFile(ctx context.Context, dataID int) (file []byte, err error) {
	action := "get file"

	file, err = sv.rp.GetDataFile(ctx, dataID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return file, nil
}
