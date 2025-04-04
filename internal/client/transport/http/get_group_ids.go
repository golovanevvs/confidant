package trhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) GetGroupIDs(ctx context.Context, accessToken string) (groupIDs map[int]struct{}, err error) {
	action := "get group IDs"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/groups", tr.addr)

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return groupIDs, fmt.Errorf(
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
		return groupIDs, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrSendRequest,
			err,
		)
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&groupIDs)
	if err != nil {
		return groupIDs, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}

	return groupIDs, nil
}
