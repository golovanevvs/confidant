package trhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) RegisterAccount(email, password string) (response *http.Response, err error) {
	action := "register account"

	endpoint := fmt.Sprintf("http://%s/register", tr.addr)

	account := model.Account{
		Email:    email,
		Password: password,
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w: %w", customerrors.ClientHTTPErr, action, customerrors.ErrEncodeJSON500, err)
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(accountJSON))
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w: %w", customerrors.ClientHTTPErr, action, customerrors.ErrCreateRequest, err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err = tr.cl.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w: %w", customerrors.ClientHTTPErr, action, customerrors.ErrSendRequest, err)
	}
	defer response.Body.Close()

	return response, nil
}
