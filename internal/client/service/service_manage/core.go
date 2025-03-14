package service_manage

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceManage) GetServerStatus(ctx context.Context) (statusResp *model.StatusResp, err error) {
	action := "get server status"

	trResponse, err := sv.tr.GetServerStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	statusResp = &model.StatusResp{
		HTTPStatusCode: trResponse.HTTPStatusCode,
		HTTPStatus:     trResponse.HTTPStatus,
	}

	return statusResp, nil
}
