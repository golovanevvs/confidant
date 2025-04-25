package service_sync

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceSync) SyncAll(ctx context.Context, accessToken string, email string) (syncResp *model.SyncResp, err error) {
	action := "sync all"

	_, err = sv.SyncGroups(ctx, accessToken, email)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataSyncResp, err := sv.SyncData(ctx, accessToken, email)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return &model.SyncResp{
		HTTPStatusCode: dataSyncResp.HTTPStatusCode,
		HTTPStatus:     dataSyncResp.HTTPStatus,
		Error:          dataSyncResp.Error,
	}, nil
}
