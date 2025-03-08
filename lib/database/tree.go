package database

type Tree struct {
	Entries []Entry
	Oid     string
}

func NewTree(entries []Entry) *Tree {
	return &Tree{Entries: entries}
}
