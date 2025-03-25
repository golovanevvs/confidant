package model

type Group struct {
	ID         int
	IDOnServer int
	AccountID  int
	Title      string
	Emails     []string
}
