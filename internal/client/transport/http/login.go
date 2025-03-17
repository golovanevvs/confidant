package trhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) Login(ctx context.Context, email, password string) (trResponse *model.TrResponse, err error) {
	//! Request
	action := "login"

	endpoint := fmt.Sprintf("http://%s/api/login", tr.addr)

	account := model.Account{
		Email:    email,
		Password: password,
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrEncodeJSON500,
			err,
		)
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(accountJSON))
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrCreateRequest,
			err,
		)
	}

	request.Header.Set("Content-Type", "application/json")

	//! Response
	response, err := tr.cl.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrSendRequest,
			err,
		)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrReadResponseBody,
			err,
		)
	}

	authHeader := response.Header.Get("Authorization")
	refreshTokenHeader := response.Header.Get("Refresh-Token")

	//! Result
	trResponse = &model.TrResponse{
		HTTPStatusCode:     response.StatusCode,
		HTTPStatus:         response.Status,
		AuthHeader:         authHeader,
		RefreshTokenHeader: refreshTokenHeader,
		ResponseBody:       responseBody,
	}

	return trResponse, nil
}
