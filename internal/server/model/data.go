package model

import "time"

type NoteDec struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       string
	Note       string
}

type NoteEnc struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       []byte `json:"desc"`
	Note       []byte `json:"note"`
}

type PassDec struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       string
	Login      string
	Pass       string
}

type PassEnc struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       []byte `json:"desc"`
	Login      []byte `json:"login"`
	Pass       []byte `json:"pass"`
}

type CardDec struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       string
	Number     string
	Date       string
	Name       string
	CVC2       string
	PIN        string
	Bank       string
}

type CardEnc struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       []byte `json:"desc"`
	Number     []byte `json:"number"`
	Date       []byte `json:"date"`
	Name       []byte `json:"name"`
	CVC2       []byte `json:"cvc2"`
	PIN        []byte `json:"pin"`
	Bank       []byte `json:"bank"`
}

type FileDec struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       string
	Filename   string
	Filesize   string
	Filedate   string
	File       []byte
}

type FileEnc struct {
	ID         int
	IDOnServer int
	GroupID    int
	Type       string
	Title      string
	Desc       []byte `json:"desc"`
	Filename   []byte `json:"filename"`
	Filesize   []byte `json:"filesize"`
	Filedate   []byte `json:"filedate"`
	File       []byte `json:"file"`
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
