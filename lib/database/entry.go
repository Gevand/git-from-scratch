package database

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
)

const (
	REGULAR_MODE    = "100644"
	EXECUTABLE_MODE = "100755"
)

type Entry struct {
	Name, Oid string
	Mode      os.FileMode
}

func NewEntry(path, oid string, mode os.FileMode) *Entry {
	return &Entry{Name: path, Oid: oid, Mode: mode}
}
func (t *Tree) ToString() string {
	return_value := ""
	sort.Slice(t.Entries, func(i, j int) bool {
		return t.Entries[i].Name < t.Entries[j].Name
	})
	for _, entry := range t.Entries {
		string_mode := REGULAR_MODE
		if entry.Mode&0111 != 0 {
			string_mode = EXECUTABLE_MODE
		}

		temp_string := fmt.Sprintf("%v %v\000", string_mode, entry.Name)
		oid_as_array := []byte(entry.Oid)
		oid_as_hexstring := hex.EncodeToString(oid_as_array[:min(20, len(oid_as_array))])
		return_value += temp_string
		return_value += oid_as_hexstring
	}
	return return_value
}
