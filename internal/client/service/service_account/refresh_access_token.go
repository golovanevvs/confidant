package service_account

import (
	"context"
	"fmt"
	"strings"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceAccount) RefreshAccessToken(ctx context.Context, refreshToken string) (refreshAccessTokenResp *model.RefreshAccessTokenResp, err error) {
	action := "refresh access token"

	var trResponse *model.TrResponse
	trResponse, err = sv.tr.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}
	if trResponse.HTTPStatusCode != 200 {
		return &model.RefreshAccessTokenResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
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

	if len(authHeader) == 0 {
		return &model.RefreshAccessTokenResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			Error: fmt.Sprintf(
				"%s: %s: %s: %s",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrAuthHeaderResp.Error(),
			),
		}, nil
	}

	if len(authHeaderSplit) != 2 {
		return &model.RefreshAccessTokenResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			Error: fmt.Sprintf(
				"%s: %s: %s: %s",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrInvalidAuthHeaderResp.Error(),
			),
		}, nil
	}

	if authHeaderSplit[0] != "Bearer" {
		return &model.RefreshAccessTokenResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			Error: fmt.Sprintf(
				"%s: %s: %s: %s",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrBearer.Error(),
			),
		}, nil
	}

	if len(authHeaderSplit[1]) == 0 {
		return &model.RefreshAccessTokenResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			Error: fmt.Sprintf(
				"%s: %s: %s: %s",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				customerrors.ErrAccessToken.Error(),
			),
		}, nil
	}

	return &model.RefreshAccessTokenResp{
		HTTPStatusCode:    trResponse.HTTPStatusCode,
		HTTPStatus:        trResponse.HTTPStatus,
		AccessTokenString: authHeaderSplit[1],
		Error:             "",
	}, nil
}
