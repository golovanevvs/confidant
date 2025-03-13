package service_account

import (
	"fmt"
	"strings"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceAccount) CreateAccount(email, password string) (registerAccountResp *model.RegisterAccountResp, err error) {
	action := "register account"

	trResponse, err := sv.tr.CreateAccount(email, password)
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

	passwordHash, err := sv.genHash(password)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientServiceErr,
			action,
			customerrors.ErrGenPasswordHash,
			err)
	}

	err = sv.rp.SaveAccount(email, passwordHash, refreshTokenHeader)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return &model.RegisterAccountResp{
		HTTPStatusCode:     trResponse.HTTPStatusCode,
		HTTPStatus:         trResponse.HTTPStatus,
		AccessTokenString:  authHeaderSplit[1],
		RefreshTokenString: refreshTokenHeader,
		Error:              "",
	}, nil
}
