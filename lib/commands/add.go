package commands

import (
	"fmt"
	"geo-git/lib"
	db "geo-git/lib/database"
	"os"
	"path"
)

func RunAdd(root_path string) error {
	git_path := path.Join(root_path, ".git")
	db_path := path.Join(git_path, "objects")
	index_path := path.Join(git_path, "index")

	workspace := lib.NewWorkSpace(root_path)
	database := db.NewDatabase(db_path)
	index := lib.NewIndex(index_path)

	path := os.Args[2]
	fmt.Println("Running add with ", path)
	data, err := workspace.ReadFile(path)
	if err != nil {
		return err
	}

	blob := db.NewBlob(data)
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	fmt.Print("Add blob", blob.Data, string(blob.Data))
	err = database.StoreBlob(blob)
	if err != nil {
		return err
	}
	err = index.Add(path, blob.Oid, stat)
	if err != nil {
		return err
	}

	_, err = index.WriteUpdates()
	return err
}
