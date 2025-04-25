package trhttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) SendFile(ctx context.Context, accessToken string, dataID int, file []byte) (err error) {
	action := "send datas"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/data_file", tr.addr)

	request, err := http.NewRequest("PUT", endpoint, bytes.NewReader(file))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrCreateRequest,
			err,
		)
	}

	request.Header.Set("X-Data-ID", strconv.Itoa(dataID))
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	//! Response
	_, err = tr.cl.Do(request)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrSendRequest,
			err,
		)
	}

	return nil
}
