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
