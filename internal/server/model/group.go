package model

type Group struct {
	ID         int      `json:"server_id"`
	IDOnClient int      `json:"client_id"`
	AccountID  int      `json:"account_id"`
	Title      string   `json:"title"`
	Emails     []string `json:"emails"`
}
