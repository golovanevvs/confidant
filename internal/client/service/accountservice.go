package service

import (
	"encoding/json"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *Service) RegisterAccount(email, password string) (registerAccountResp *model.RegisterAccountResp, err error) {
	action := "register account"

	trResponse, err := sv.tr.RegisterAccount(email, password)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", customerrors.ClientHTTPErr, action, err)
	}

	if trResponse.HTTPStatusCode == 200 {
		var accountRegisterResp model.AccountRegisterResp
		err = json.Unmarshal(trResponse.ResponseBody, &accountRegisterResp)
		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w: %w", customerrors.ClientHTTPErr, action, customerrors.ErrDecodeJSON400, err)
		}
		registerAccountResp = &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      accountRegisterResp.AccountID,
			ServerError:    "",
		}
		return registerAccountResp, nil
	} else {
		registerAccountResp = &model.RegisterAccountResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			AccountID:      "",
			ServerError:    string(trResponse.ResponseBody),
		}
		return registerAccountResp, nil
	}
}
