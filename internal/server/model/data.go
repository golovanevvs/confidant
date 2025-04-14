package model

import "time"

type NoteEnc struct {
	ID   int
	Desc []byte `json:"desc"`
	Note []byte `json:"note"`
}

type PassEnc struct {
	ID    int
	Desc  []byte `json:"desc"`
	Login []byte `json:"login"`
	Pass  []byte `json:"pass"`
}

type CardEnc struct {
	ID     int
	Desc   []byte `json:"desc"`
	Number []byte `json:"number"`
	Date   []byte `json:"date"`
	Name   []byte `json:"name"`
	CVC2   []byte `json:"cvc2"`
	PIN    []byte `json:"pin"`
	Bank   []byte `json:"bank"`
}

type FileEnc struct {
	ID       int
	Desc     []byte `json:"desc"`
	Filename []byte `json:"filename"`
	Filesize []byte `json:"filesize"`
	Filedate []byte `json:"filedate"`
	File     []byte `json:"-"`
}

type Data struct {
	ID         int       `json:"id_on_server"`
	IDOnClient int       `json:"id_on_client"`
	GroupID    int       `json:"group_id"`
	DataType   string    `json:"data_type"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	Note       NoteEnc   `json:"note"`
	Pass       PassEnc   `json:"pass"`
	Card       CardEnc   `json:"card"`
	File       FileEnc   `json:"file"`
}
