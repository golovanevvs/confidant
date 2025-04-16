package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataFile(ctx context.Context, dataID int) (idOnServer int, file []byte, err error) {
	action := "get file"

	idOnServer, file, err = sv.rp.GetDataFile(ctx, dataID)
	if err != nil {
		return -1, nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return idOnServer, file, nil
}
