package database

import (
	"encoding/hex"
	"fmt"
	"sort"
)

type Entry struct {
	Name, Oid string
}

func NewEntry(path, oid string) *Entry {
	return &Entry{Name: path, Oid: oid}
}
func (t *Tree) ToString() string {
	return_value := ""
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Name < t.Entries[j].Name
	})
	for _, entry := range t.Entries {
		temp_string := fmt.Sprintf("%v %v\000", t.mode, entry.Name)
		oid_as_array := []byte(entry.Oid)
		oid_as_hexstring := hex.EncodeToString(oid_as_array[:min(20, len(oid_as_array))])
		return_value += temp_string
		return_value += oid_as_hexstring
	}
	return return_value
}
