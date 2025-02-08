package commands

import (
	"fmt"
	"os"
	"path"
)

var dirs = []string{"objects", "refs"}

func RunInit(root_path, git_path string) error {
	err := os.Mkdir(git_path, 0777)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		err := os.Mkdir(path.Join(git_path, dir), 0777)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Initialized empty repository in %s\r\n", git_path)
	return nil
}
