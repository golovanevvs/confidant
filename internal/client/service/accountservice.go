package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *Service) RegisterAccount(email, password string) (registerAccountResp *model.RegisterAccountResp, err error) {
	action := "register account service"

	trResponse, err := sv.tr.RegisterAccount(email, password)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", customerrors.ClientServiceErr, action, err)
	}

	if trResponse.HTTPStatusCode != 200 {
		return &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      "",
			Error: fmt.Sprintf(
				"%s: %s: %s: %s",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				string(trResponse.ResponseBody),
			),
		}, nil
	}

	authHeader := trResponse.AuthHeader
	authHeaderSplit := strings.Split(authHeader, " ")
	refreshTokenHeader := trResponse.RefreshTokenHeader

	if authHeader == "" {
		return &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      "",
			Error: fmt.Sprintf(
				"%s: %s: %s",
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrAuthHeader.Error(),
			),
		}, nil
	}

	if authHeader == "" {
		return &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      "",
			Error: fmt.Sprintf(
				"%s: %s: %s",
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrInvalidAuthHeader.Error(),
			),
		}, nil
	}

	if authHeaderSplit[0] != "Bearer" {
		return &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      "",
			Error: fmt.Sprintf(
				"%s: %s: %s",
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrBearer.Error(),
			),
		}, nil
	}

	if authHeaderSplit[1] == "" {
		return &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      "",
			Error: fmt.Sprintf(
				"%s: %s: %s",
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrAccessToken.Error(),
			),
		}, nil
	}

	if refreshTokenHeader == "" {
		return &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      "",
			Error: fmt.Sprintf(
				"%s: %s: %s",
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrRefreshToken.Error(),
			),
		}, nil
	}

	var accountRegisterResp model.AccountRegisterResp
	err = json.Unmarshal(trResponse.ResponseBody, &accountRegisterResp)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}

	//! Сохранить access токен в SQLite
	//! Сохранить refresh токен в SQLite

	return &model.RegisterAccountResp{
		HTTPStatusCode: trResponse.HTTPStatusCode,
		HTTPStatus:     trResponse.HTTPStatus,
		AccountID:      accountRegisterResp.AccountID,
		Error:          "",
	}, nil
}
