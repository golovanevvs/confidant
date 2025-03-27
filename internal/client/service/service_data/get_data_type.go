package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataIDAndType(ctx context.Context, groupID int, dataTitle string) (dataID int, dataType string, err error) {
	action := "get data ID and data type"

	dataTitleEnc, err := sv.ss.Encrypt([]byte(dataTitle))
	if err != nil {
		return -1, "", fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataID, dataType, err = sv.rp.GetDataIDAndType(ctx, groupID, dataTitleEnc)
	if err != nil {
		return -1, "", fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return dataID, dataType, nil
}
