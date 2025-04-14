package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) SaveDatas(ctx context.Context, datas []model.Data) (err error) {
	action := "save datas"

	err = sv.rp.SaveDatas(ctx, datas)
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
