package commands

import (
	"fmt"
	"geo-git/lib"
	"os"
	"path"
)

var dirs = []string{"objects", "refs"}

func RunInit(repo *lib.Respository, cmd *Command) error {
	err := os.Mkdir(repo.GitPath, 0777)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		err := os.Mkdir(path.Join(repo.GitPath, dir), 0777)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Initialized empty repository in %s\r\n", repo.GitPath)
	return nil
}
