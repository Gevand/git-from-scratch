package database

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Blob struct {
	Data []byte
	Oid  string
	Type string
}

func NewBlob(data string) *Blob {
	return &Blob{Data: []byte(data), Type: "blob"}
}

func (b *Blob) GetContent() string {
	content := fmt.Sprintf("%s %v\000%s", b.Type, len(b.Data), b.Data)
	return content
}

func (b *Blob) HashObject() string {
	h := sha1.New()
	h.Write([]byte(b.GetContent()))
	return hex.EncodeToString(h.Sum(nil))
}
