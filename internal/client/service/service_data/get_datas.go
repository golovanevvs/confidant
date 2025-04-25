package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDatas(ctx context.Context, dataIDs []int) (datas []model.Data, err error) {
	action := "get datas"

	datas, err = sv.rp.GetDatasByIDs(ctx, dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return datas, nil
}
