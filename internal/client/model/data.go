package model

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
