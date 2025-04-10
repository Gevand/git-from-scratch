package commands

import (
	"fmt"
	"geo-git/lib"
	db "geo-git/lib/database"
	"os"
	"path"
)

func RunCommit(root_path string, author *db.Author, message string) error {

	git_path := path.Join(root_path, ".git")
	db_path := path.Join(git_path, "objects")

	workspace := lib.NewWorkSpace(root_path)
	database := db.NewDatabase(db_path)
	refs := lib.NewRefs(git_path)
	parent, err := refs.ReadHead()
	if err != nil {
		return err
	}

	files, err := workspace.ListFiles(root_path)
	if err != nil {
		return err
	}
	fmt.Println("Commit files", files)

	tree_entries := []*db.Entry{}
	for _, file := range files {
		data, err := workspace.ReadFile(file)
		if err != nil {
			return err
		}
		blob := db.NewBlob(data)
		err = database.StoreBlob(blob)
		if err != nil {
			return err
		}
		stats, err := os.Stat(file)
		if err != nil {
			return err
		}
		mode := stats.Mode().Perm()
		entry := db.NewEntry(file, blob.Oid, mode)
		tree_entries = append(tree_entries, entry)
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
