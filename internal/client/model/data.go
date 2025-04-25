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
	Desc       []byte
	Note       []byte
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
	Desc       []byte
	Login      []byte
	Pass       []byte
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
	Desc       []byte
	Number     []byte
	Date       []byte
	Name       []byte
	CVC2       []byte
	PIN        []byte
	Bank       []byte
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
	Desc       []byte
	Filename   []byte
	Filesize   []byte
	Filedate   []byte
	File       []byte
}

type Data struct {
	ID         int
	IDOnServer int
	GroupID    int
	DataType   string
	Title      string
	CreatedAt  time.Time
	Note       NoteEnc
	Pass       PassEnc
	Card       CardEnc
	File       FileEnc
}
