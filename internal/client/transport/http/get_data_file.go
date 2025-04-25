package trhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) GetDataFile(ctx context.Context, accessToken string, dataID int) (fileFromServer []byte, err error) {
	action := "get file"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/data_file", tr.addr)

	dataIDJSON, err := json.Marshal(dataID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrEncodeJSON500,
			err,
		)
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(dataIDJSON))
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

	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientHTTPErr,
			action,
			err,
		)
	}

	fileFromServer = buf.Bytes()

	return fileFromServer, nil
}
