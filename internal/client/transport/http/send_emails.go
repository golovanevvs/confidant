package trhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) SendEmails(ctx context.Context, accessToken string, mapGroupIDEmails map[int][]string) (err error) {
	action := "send groups"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/emails", tr.addr)

	mapGroupIDEmailsJSON, err := json.Marshal(mapGroupIDEmails)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrEncodeJSON500,
			err,
		)
	}

	request, err := http.NewRequest("PATCH", endpoint, bytes.NewReader(mapGroupIDEmailsJSON))
	if err != nil {
		return fmt.Errorf(
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
