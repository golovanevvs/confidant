package trhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) RefreshAccessToken(ctx context.Context, refreshToken string) (trResponse *model.TrResponse, err error) {
	action := "refresh access token"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/refresh_access", tr.addr)

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(refreshToken)))
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrCreateRequest,
			err,
		)
	}

	request.Header.Set("Content-Type", "text/plain")

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

	//! Result
	trResponse = &model.TrResponse{
		HTTPStatusCode: response.StatusCode,
		HTTPStatus:     response.Status,
		AuthHeader:     authHeader,
		ResponseBody:   responseBody,
	}

	return trResponse, nil
}
