package commands

import (
	"fmt"
	"geo-git/lib"
	db "geo-git/lib/database"
	"geo-git/lib/utils"
	"os"
	"path"
)

func RunCommit(root_path string, author *db.Author, message string) error {

	git_path := path.Join(root_path, ".git")
	db_path := path.Join(git_path, "objects")
	index_path := path.Join(git_path, "index")

	index := lib.NewIndex(index_path)
	err := index.Load()
	if err != nil {
		return err
	}
	database := db.NewDatabase(db_path)
	refs := lib.NewRefs(git_path)

	parent, err := refs.ReadHead()
	if err != nil {
		return err
	}

	tree_entries := []*db.Entry{}
	for _, index_entry := range utils.SliceFromMap(index.Entries) {
		tree_entries = append(tree_entries, db.NewEntry(index_entry.Path, index_entry.Oid, os.FileMode(index_entry.Mode)))
	}
	root := db.NewTree("")
	root.BuildTree(tree_entries)
	err = root.Traverse(func(t *db.Tree) error {
		err := database.StoreTree(t)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("tree:", root.Oid)

	commit := db.NewCommit(parent, root.Oid, *author, message)
	err = database.StoreCommit(commit)
	if err != nil {
		return err
	}
	err = refs.UpdateHead(commit.Oid)
	if err != nil {
		return err
	}

	head_file, err := os.OpenFile(path.Join(git_path, "HEAD"), os.O_WRONLY|os.O_CREATE, 0777)
	defer head_file.Close()
	head_file.Write([]byte(commit.Oid))

	fmt.Println("commit:", commit.Oid, message)
	return nil
}
