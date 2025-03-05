package model

type TrResponse struct {
	HTTPStatusCode int
	HTTPStatus     string
	ResponseBody   []byte
}

type RegisterAccountResp struct {
	HTTPStatusCode int
	HTTPStatus     string
	AccountID      string
	ServerError    string
}

type AccountRegisterResp struct {
	AccountID string `json:"account_id"`
}
