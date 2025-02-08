package commands

import (
	"fmt"
	"geo-git/lib"
	"path"
)

func RunCommit(root_path string) error {

	git_path := path.Join(root_path, ".git")
	db_path := path.Join(git_path, "objects")

	workspace := lib.NewWorkSpace(root_path)
	database := lib.NewDatabase(db_path)

	files, err := workspace.ListFiles()
	if err != nil {
		return err
	}
	fmt.Println("Commit files", files)

	tree_entries := []lib.Entry{}
	for _, file := range files {
		data, err := workspace.ReadFile(file)
		if err != nil {
			return err
		}
		blob := lib.NewBlob(data)
		err = database.StoreBlob(blob)
		if err != nil {
			return err
		}
		tree_entries = append(tree_entries, *lib.NewEntry(file, blob.Oid))
	}
	tree := lib.NewTree(tree_entries)
	err = database.StoreTree(tree)
	if err != nil {
		return err
	}
	return nil
}
