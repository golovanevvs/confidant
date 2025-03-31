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
