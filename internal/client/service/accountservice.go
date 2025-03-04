package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *Service) RegisterAccount(email, password string) (result model.RegisterAccount, err error) {
	response, err := sv.tr.RegisterAccount(email, password)

	if response.StatusCode == http.StatusOK {

		respBody := response.Body

		var responseData struct {
			AccountID string `json:"accountid"`
			Token     string `json:"token"`
		}
		dec := json.NewDecoder(respBody)
		err = dec.Decode(&responseData)
		if err != nil {
			return nil, fmt.Errorf("%s: %s: %w: %w", customerrors.ClientHTTPErr, action, customerrors.ErrDecodeJSON400, err)
		}

		return accountID, nil
	} else {
		responseError, err := io.ReadAll(response.Body)
		if err != nil {
			return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.ClientHTTPErr, action, customerrors.ErrReadResponseBody, err)
		}
		return -1, fmt.Errorf("%s: %s: %s[yellow]status: %d", customerrors.ClientHTTPErr, action, string(responseError), response.StatusCode)
	}
	return accountID, err
}
