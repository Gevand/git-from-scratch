package database

import (
	"os"
	"path/filepath"
	"slices"
)

const (
	REGULAR_MODE    = "100644"
	EXECUTABLE_MODE = "100755"
	DIRECTORY_MODE  = "040000"
)

type Entry struct {
	Name, Oid string
	Mode      os.FileMode
}

func NewEntry(path, oid string, mode os.FileMode) *Entry {
	return &Entry{Name: path, Oid: oid, Mode: mode}
}

func (e *Entry) ParentDirectories() []string {
	subPath := e.Name
	var result []string
	for {
		subPath = filepath.Clean(subPath)

		dir, last := filepath.Split(subPath)
		if last == "" {
			if dir != "" {
				result = append(result, dir)
			}
			break
		}
		result = append(result, last)

		if dir == "" {
			break
		}
		subPath = dir
	}

	slices.Reverse(result)
	return result[:len(result)-1]
}
