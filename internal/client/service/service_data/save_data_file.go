package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) SaveDataFile(ctx context.Context, dataID int, file []byte) (err error) {
	action := "save file"

	err = sv.rp.SaveDataFile(ctx, dataID, file)
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
