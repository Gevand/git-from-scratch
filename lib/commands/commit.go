package commands

import (
	"fmt"
	"geo-git/lib"
	"os"
	"path"
)

func RunCommit(root_path string) {

	git_path := path.Join(root_path, ".git")
	db_path := path.Join(git_path, "objects")

	workspace := lib.NewWorkSpace(root_path)
	database := lib.NewDatabase(db_path)

	files, err := workspace.ListFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "commit - fatal: %v\r\n", err)
		return
	}
	fmt.Println("Commit files", files)
	for _, file := range files {
		data, err := workspace.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit - fatal: %v\r\n", err)
			return
		}
		blob := lib.NewBlob(data)
		err = database.Store(blob)
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit - fatal: %v\r\n", err)
			return
		}
	}
}
