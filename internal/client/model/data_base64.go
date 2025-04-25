package model

import "time"

type NoteBase64 struct {
	Desc string `json:"desc"`
	Note string `json:"note"`
}

type PassBase64 struct {
	Desc  string `json:"desc"`
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type CardBase64 struct {
	Desc   string `json:"desc"`
	Number string `json:"number"`
	Date   string `json:"date"`
	Name   string `json:"name"`
	CVC2   string `json:"cvc2"`
	PIN    string `json:"pin"`
	Bank   string `json:"bank"`
}

type FileBase64 struct {
	Desc     string `json:"desc"`
	Filename string `json:"filename"`
	Filesize string `json:"filesize"`
	Filedate string `json:"filedate"`
	File     []byte `json:"-"`
}

type DataBase64 struct {
	ID         int        `json:"id_on_client"`
	IDOnServer int        `json:"id_on_server"`
	GroupID    int        `json:"group_id"`
	DataType   string     `json:"data_type"`
	Title      string     `json:"title"`
	CreatedAt  time.Time  `json:"created_at"`
	Note       NoteBase64 `json:"notebase64"`
	Pass       PassBase64 `json:"passbase64"`
	Card       CardBase64 `json:"cardbase64"`
	File       FileBase64 `json:"filebase64"`
}
