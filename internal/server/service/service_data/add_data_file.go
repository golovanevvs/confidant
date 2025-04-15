package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) AddDataFile(ctx context.Context, dataID int, file []byte) (err error) {
	action := "add file"

	err = sv.rp.SaveDataFile(ctx, dataID, file)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.DataServiceErr,
			action,
			err,
		)
	}

	return nil
}
