package database

type Tree struct {
	Entries []Entry
	Oid     string
	mode    string
}

func NewTree(entries []Entry) *Tree {
	return &Tree{Entries: entries, mode: "100644"}
}
