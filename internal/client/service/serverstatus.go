package service

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *Service) GetServerStatus() (statusResp *model.StatusResp, err error) {
	action := "get server status service"

	trResponse, err := sv.tr.ServerStatus()
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
