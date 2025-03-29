package commands

import (
	"fmt"
	"geo-git/lib"
	db "geo-git/lib/database"
	"os"
	"path"
	"path/filepath"
)

func RunAdd(root_path string) error {
	git_path := path.Join(root_path, ".git")
	db_path := path.Join(git_path, "objects")
	index_path := path.Join(git_path, "index")

	workspace := lib.NewWorkSpace(root_path)
	database := db.NewDatabase(db_path)
	index := lib.NewIndex(index_path)

	for i := 2; i < len(os.Args); i++ {
		path_from_arg := os.Args[i]
		if !filepath.IsAbs(path_from_arg) {
			absolute, err := filepath.Abs(path_from_arg)
			if err != nil {
				return err
			}
			path_from_arg = absolute
		}

		//expand the path if its a folder
		all_paths, err := workspace.ListFiles(path_from_arg)
		if err != nil {
			return err
		}

		for _, path := range all_paths {
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
		}
	}

	_, err := index.WriteUpdates()
	return err
}
