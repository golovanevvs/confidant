package trhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type trHTTP struct {
	addr string
}

func New() *trHTTP {
	return &trHTTP{}
}

func (t *trHTTP) RegisterAccount(account model.Account) {
	endpoint := fmt.Sprintf("http://:8080/register")

	accountJSON, err := json.Marshal(account)
	if err != nil {
		fmt.Printf("Ошибка кодирования в JSON: %s\n", err.Error())
		return
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(accountJSON))
	if err != nil {
		fmt.Printf("Ошибка формирования запроса: %s\n", err.Error())
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Ошибка отправки запроса: %s\n", err.Error())
		return
	}
	defer response.Body.Close()

	respBody := response.Body

	dec := json.NewDecoder(respBody)
	var responseData struct {
		Email     string `json:"email"`
		AccountID string `json:"accountid"`
		Token     string `json:"token"`
	}

	err = dec.Decode(&responseData)
	if err != nil {
		fmt.Printf("Ошибка декодирования JSON: %s\n", err.Error())
		return
	}

	fmt.Printf("Response: %v\n", responseData)
}
