package lib

import (
	"geo-git/lib/utils"
	"os"
)

var ignore = []string{".", "..", ".git", "test_script.sh", "git-from-scratch"}

type Workspace struct {
	pathname string
}

func (w *Workspace) ReadFile(file string) (string, error) {
	str, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return (string)(str), err
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
			continue
		}
		result = append(result, file.Name())
	}
	return result, nil
}
