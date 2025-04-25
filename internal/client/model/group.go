package model

type Group struct {
	ID         int      `json:"client_id"`
	IDOnServer int      `json:"server_id"`
	AccountID  int      `json:"account_id"`
	Title      string   `json:"title"`
	Emails     []string `json:"emails"`
}
