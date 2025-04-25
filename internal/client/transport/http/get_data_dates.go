package trhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) GetDataDates(ctx context.Context, accessToken string, dataIDs []int) (dataDatesFromServer map[int]time.Time, err error) {
	action := "get data dates"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/data_dates", tr.addr)

	dataIDsJSON, err := json.Marshal(dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrEncodeJSON500,
			err,
		)
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(dataIDsJSON))
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

	err = json.NewDecoder(response.Body).Decode(&dataDatesFromServer)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}

	return dataDatesFromServer, nil
}
