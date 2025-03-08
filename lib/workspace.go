package lib

import (
	"fmt"
	"geo-git/lib/utils"
	"os"
	"path/filepath"
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

func (w *Workspace) ListFiles(path string) ([]string, error) {
	result := []string{}
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if utils.Contains[string](ignore, file.Name()) {
			continue
		}
		if file.IsDir() {
			temp_path := filepath.Join(path, file.Name())
			temp_results, err := w.ListFiles(temp_path)
			if err != nil {
				return nil, err
			}
			result = append(result, temp_results...)
		} else {
			full_path := filepath.Join(path, file.Name())
			relative_path, err := filepath.Rel(w.pathname, full_path)
			if err != nil {
				return nil, err
			}
			result = append(result, relative_path)
		}
	}
	fmt.Println("results: ", result)
	return result, nil
}
