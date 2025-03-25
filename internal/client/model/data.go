package model

type NoteDec struct {
	ID         int
	IDOnServer int
	GroupsID   int
	Title      string
	Desc       string
	Note       string
}

type NoteEnc struct {
	ID         int
	IDOnServer int
	GroupsID   int
	Title      []byte
	Desc       []byte
	Note       []byte
}
