package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDataTitles(ctx context.Context, accountID int, groupID int) (dataTitles []string, err error) {
	action := "get data titles"

	dataTitlesB, err := sv.rp.GetDataTitles(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	for _, data := range dataTitlesB {
		dataTitles = append(dataTitles, string(data))
	}

	// for _, dataTitleEnc := range dataTitlesEnc {
	// 	dataTitleDec, err := sv.ss.Decrypt(dataTitleEnc)
	// 	if err != nil {
	// 		return nil, fmt.Errorf(
	// 			"%s: %s: %w",
	// 			customerrors.ClientServiceErr,
	// 			action,
	// 			err,
	// 		)
	// 	}
	// 	dataTitles = append(dataTitles, string(dataTitleDec))
	// }

	return dataTitles, nil
}
