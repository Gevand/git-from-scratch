package commands

import (
	"fmt"
	"geo-git/lib"
	"sort"
)

func RunStatus(repo *lib.Respository, cmd *Command) error {
	err := repo.Index.Load()
	if err != nil {
		return err
	}
	files, err := repo.Workspace.ListFiles(repo.Workspace.Pathname)
	sort.Strings(files)
	if err != nil {
		return err
	}
	for _, file := range files {
		if repo.Index.IsEntryTracked(file) {
			fmt.Printf("%s\r\n", file)
		} else {
			fmt.Printf("?? %s\r\n", file)
		}
	}
	return nil
}
