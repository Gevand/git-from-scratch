package commands

import (
	"geo-git/lib"
	db "geo-git/lib/database"
	"os"
	"path/filepath"
)

func RunAdd(repo *lib.Respository, cmd *Command) error {
	repo.Index.LoadForUpdate()
	for _, arg := range cmd.Args {
		path_from_arg := arg
		if !filepath.IsAbs(path_from_arg) {
			absolute, err := filepath.Abs(path_from_arg)
			if err != nil {
				return err
			}
			path_from_arg = absolute
		}

		//expand the path if its a folder
		all_paths, err := repo.Workspace.ListFiles(path_from_arg)
		if err != nil {
			return err
		}

		for _, path := range all_paths {
			data, err := repo.Workspace.ReadFile(path)
			if err != nil {
				return err
			}

			blob := db.NewBlob(data)
			stat, err := os.Stat(path)
			if err != nil {
				return err
			}
			err = repo.Database.StoreBlob(blob)
			if err != nil {
				return err
			}
			err = repo.Index.Add(path, blob.Oid, stat)
			if err != nil {
				return err
			}
		}
	}

	_, err := repo.Index.WriteUpdates()
	return err
}
