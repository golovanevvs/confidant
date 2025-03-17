package model

type TrResponse struct {
	HTTPStatusCode     int
	HTTPStatus         string
	AuthHeader         string
	RefreshTokenHeader string
	ResponseBody       []byte
}

type AccountResp struct {
	HTTPStatusCode     int
	HTTPStatus         string
	AccountID          int
	AccessTokenString  string
	RefreshTokenString string
	Error              string
}

// type AccountRegisterResp struct {
// 	AccountID string `json:"account_id"`
// }

type StatusResp struct {
	HTTPStatusCode int
	HTTPStatus     string
}
