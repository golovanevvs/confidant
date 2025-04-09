package trhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) GetGroupIDs(ctx context.Context, accessToken string) (trResponse *model.GroupSyncResp, err error) {
	action := "get group IDs"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/groupids", tr.addr)

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrCreateRequest,
			err,
		)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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

	return &model.GroupSyncResp{
		HTTPStatusCode: response.StatusCode,
		HTTPStatus:     response.Status,
		ResponseBody:   responseBody,
	}, nil
}
