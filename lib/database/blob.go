package database

type Blob struct {
	Data []byte
	Oid  string
	Type string
}

func NewBlob(data string) *Blob {
	return &Blob{Data: []byte(data), Type: "blob"}
}
