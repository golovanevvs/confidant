package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataTitles(ctx context.Context, accountID int, titleGroup string) (dataTitles []string, err error) {
	action := "get data titles"

	groupID, err := sv.rp.GetGroupID(ctx, accountID, titleGroup)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataTitlesEnc, err := sv.rp.GetDataTitles(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	for _, dataTitleEnc := range dataTitlesEnc {
		dataTitleDec, err := sv.ss.Decrypt(dataTitleEnc)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w",
				customerrors.ClientServiceErr,
				action,
				err,
			)
		}
		dataTitles = append(dataTitles, string(dataTitleDec))
	}

	return dataTitles, nil
}
