package lib

import (
	"geo-git/lib/utils"
	"os"
)

// test
var ignore = []string{".", "..", ".git"}

type Workspace struct {
	pathname string
}

func NewWorkSpace(pathname string) *Workspace {
	return &Workspace{pathname: pathname}
}

func (w *Workspace) ListFiles() ([]string, error) {
	result := []string{}
	files, err := os.ReadDir(w.pathname)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if utils.Contains[string](ignore, file.Name()) {
			break
		}
		result = append(result, file.Name())
	}
	return result, nil
}
