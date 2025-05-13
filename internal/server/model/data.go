package model

import "time"

type NoteEnc struct {
	// ID     int
	DataID int    `db:"data_id"`
	Desc   []byte `json:"desc" db:"descr"`
	Note   []byte `json:"note" db:"note"`
}

type PassEnc struct {
	// ID     int
	DataID int    `db:"data_id"`
	Desc   []byte `json:"desc" db:"descr"`
	Login  []byte `json:"login" db:"login"`
	Pass   []byte `json:"pass" db:"pass"`
}

type CardEnc struct {
	// ID     int
	DataID int    `db:"data_id"`
	Desc   []byte `json:"desc" db:"descr"`
	Number []byte `json:"number" db:"number"`
	Date   []byte `json:"date" db:"date"`
	Name   []byte `json:"name" db:"name"`
	CVC2   []byte `json:"cvc2" db:"cvc2"`
	PIN    []byte `json:"pin" db:"pin"`
	Bank   []byte `json:"bank" db:"bank"`
}

type FileEnc struct {
	// ID       int
	DataID   int    `db:"data_id"`
	Desc     []byte `json:"desc" db:"descr"`
	Filename []byte `json:"filename" db:"filename"`
	Filesize []byte `json:"filesize" db:"filesize"`
	Filedate []byte `json:"filedate" db:"filedate"`
	File     []byte `json:"-"`
}

type Data struct {
	ID         int       `json:"id_on_server" db:"id"`
	IDOnClient int       `json:"id_on_client"`
	GroupID    int       `json:"group_id" db:"group_id"`
	DataType   string    `json:"data_type" db:"data_type"`
	Title      string    `json:"title" db:"title"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	Note       NoteEnc   `json:"note"`
	Pass       PassEnc   `json:"pass"`
	Card       CardEnc   `json:"card"`
	File       FileEnc   `json:"file"`
}

type DataResults struct {
	ID         int `db:"id"`
	IDOnClient int `db:"id_on_client"`
}
