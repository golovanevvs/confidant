package trhttp

import (
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) ServerStatus() (statusResp *model.TrResponse, err error) {
	//! Request
	action := "get server status"

	endpoint := fmt.Sprintf("http://%s/api/status", tr.addr)

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientHTTPErr,
			action,
			err,
		)
	}

	//! Response
	response, err := tr.cl.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientHTTPErr,
			action,
			err,
		)
	}
	defer response.Body.Close()

	//! Result
	statusResp = &model.TrResponse{
		HTTPStatusCode: response.StatusCode,
		HTTPStatus:     response.Status,
	}

	return statusResp, nil
}
