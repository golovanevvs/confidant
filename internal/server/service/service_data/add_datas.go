package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (sv *ServiceData) AddDatas(ctx context.Context, datas []model.Data) (dataIDs map[int]int, err error) {
	action := "add datas"

	dataIDs, err = sv.rp.AddDatas(ctx, datas)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.DataServiceErr,
			action,
			err,
		)
	}

	return dataIDs, nil
}
