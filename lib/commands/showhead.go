package commands

import (
	"fmt"
	"geo-git/lib"
	db "geo-git/lib/database"
	"path"
)

func RunShowHead(repo *lib.Respository) error {
	fmt.Println("Running head")
	oid, err := repo.Refs.ReadHead()
	if err != nil {
		return err
	}
	fmt.Println("HEAD IS AT", oid)
	repo.Database.Load(oid)
	blob_commit := repo.Database.Objects[oid]
	commit, err := db.ParseCommitFromBlob(blob_commit)
	if err != nil {
		return err
	}
	err = showTree(repo, commit.Tree_Oid, "")
	return err
}

func showTree(repo *lib.Respository, oid string, prefix string) error {
	if oid == "" {
		return nil
	}
	repo.Database.Load(oid)
	blob_tree := repo.Database.Objects[oid]
	tree, err := db.ParseTreeFromBlob(blob_tree)
	if err != nil {
		return err
	}

	//sort the keys so they appear in a nice order
	keys := []string{}
	for key := range tree.Entries {
		keys = append(keys, key)
	}

	for _, key := range keys {
		entry := tree.Entries[key]
		path := path.Join(prefix, key)
		switch temp_entry := entry.(type) {
		case *db.Tree:
			err := showTree(repo, temp_entry.Oid, path)
			if err != nil {
				return err
			}
		case *db.Entry:
			mode := temp_entry.Mode.String()
			fmt.Println(mode, temp_entry.Oid, path)
		}
	}
	return nil
}
