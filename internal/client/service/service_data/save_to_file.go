package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) SaveToFile(ctx context.Context, dataID int, filepath string) (err error) {
	action := "save to file"

	dataEnc, err := sv.rp.GetFileForSave(ctx, dataID)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	err = sv.ss.DecryptFile(dataEnc, filepath)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return nil
}
