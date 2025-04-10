package model

type TrResponse struct {
	HTTPStatusCode     int
	HTTPStatus         string
	AccountID          int
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

type RefreshAccessTokenResp struct {
	HTTPStatusCode    int
	HTTPStatus        string
	AccessTokenString string
	Error             string
}

type GroupSyncResp struct {
	HTTPStatusCode int
	HTTPStatus     string
	Error          string
	ResponseBody   []byte
	GroupIDs       []int
}

type SyncResp struct {
	HTTPStatusCode int
	HTTPStatus     string
	Error          string
}

// type AccountRegisterResp struct {
// 	AccountID string `json:"account_id"`
// }

type StatusResp struct {
	HTTPStatusCode int
	HTTPStatus     string
}

type AccID struct {
	AccountID int `json:"account_id"`
}
