package service_account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceAccount) Login(ctx context.Context, email, password string) (registerAccountResp *model.AccountResp, err error) {
	action := "login"

	passwordHash, err := sv.genHash(password)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			customerrors.ErrGenPasswordHash,
			err,
		)
	}

	var refreshToken string
	accountID, err := sv.rp.LoadAccountID(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, customerrors.ErrDBEmailNotFound401) {
			// login on the server
			var trResponse *model.TrResponse
			trResponse, err = sv.tr.Login(ctx, email, password)
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
				return &model.AccountResp{
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
			refreshTokenHeader := trResponse.RefreshTokenHeader

			if len(authHeader) == 0 {
				return &model.AccountResp{
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
				return &model.AccountResp{
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
				return &model.AccountResp{
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
				return &model.AccountResp{
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

			if len(refreshTokenHeader) == 0 {
				return &model.AccountResp{
					HTTPStatusCode: trResponse.HTTPStatusCode,
					HTTPStatus:     trResponse.HTTPStatus,
					Error: fmt.Sprintf(
						"%s: %s: %s: %s",
						customerrors.ClientMsg,
						customerrors.ClientServiceErr,
						action,
						customerrors.ErrRefreshToken.Error(),
					),
				}, nil
			}

			// saving the refresh token in a local DB

			var account model.Account
			err = json.Unmarshal(trResponse.ResponseBody, &account)
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

			err = sv.rp.SaveActiveAccount(ctx, account.ID)
			if err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %s: %w",
					customerrors.ClientMsg,
					customerrors.ClientServiceErr,
					action,
					err,
				)
			}

			// saving the account in a local DB
			err = sv.rp.SaveAccount(ctx, account.ID, email, passwordHash, refreshTokenHeader)
			if err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %s:%w",
					customerrors.ClientMsg,
					customerrors.ClientServiceErr,
					action,
					err,
				)
			}

			return &model.AccountResp{
				HTTPStatusCode:     trResponse.HTTPStatusCode,
				HTTPStatus:         trResponse.HTTPStatus,
				AccessTokenString:  authHeaderSplit[1],
				RefreshTokenString: refreshTokenHeader,
				AccountID:          account.ID,
				Error:              "",
			}, nil
		} else {
			return nil, fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ClientMsg,
				customerrors.AccountServiceErr,
				action,
				err,
			)
		}
	} else {
		err = sv.rp.SaveActiveAccount(ctx, accountID)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				err,
			)
		}

		_, refreshToken, err = sv.rp.LoadActiveAccount(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				err,
			)
		}
	}

	return &model.AccountResp{
		AccountID:          accountID,
		RefreshTokenString: refreshToken,
		Error:              "",
	}, nil
}
