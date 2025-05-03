package lib

import (
	"geo-git/lib/utils"
	"os"
	"path/filepath"
)

var ignore = []string{".", "..", ".git", "test_script.sh", "geo-git"}

type Workspace struct {
	Pathname string
}

func (w *Workspace) ReadFile(file string) (string, error) {
	str, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return (string)(str), err
}

func NewWorkSpace(pathname string) *Workspace {
	return &Workspace{Pathname: pathname}
}

func (w *Workspace) ListFiles(path string) ([]string, error) {
	result := []string{}
	fileinfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	//single file
	if !fileinfo.IsDir() {
		relative_path, err := filepath.Rel(w.Pathname, path)
		if err != nil {
			return nil, err
		}
		result = append(result, relative_path)
		return result, nil
	}

	//directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if utils.Contains[string](ignore, file.Name()) {
			continue
		}

		// Could possibly recurse regardless of file type
		if file.IsDir() {
			temp_path := filepath.Join(path, file.Name())
			temp_results, err := w.ListFiles(temp_path)
			if err != nil {
				return nil, err
			}
			result = append(result, temp_results...)
		} else {
			full_path := filepath.Join(path, file.Name())
			relative_path, err := filepath.Rel(w.Pathname, full_path)
			if err != nil {
				return nil, err
			}
			result = append(result, relative_path)
		}
	}
	return result, nil
}
