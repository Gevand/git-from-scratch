package commands

import (
	"bufio"
	"errors"
	"fmt"
	"geo-git/lib"
	db "geo-git/lib/database"
	"geo-git/lib/utils"
	"os"
	"path"
	"time"
)

func RunCommit(repo *lib.Respository, cmd *Command) error {
	err := repo.Index.Load()
	if err != nil {
		return err
	}

	parent, err := repo.Refs.ReadHead()
	if err != nil {
		return err
	}

	name, name_ok := cmd.EnvVars["GIT_AUTHOR_NAME"]
	email, email_ok := cmd.EnvVars["GIT_AUTHOR_EMAIL"]
	if !name_ok || !email_ok || name == "" || email == "" {
		return errors.New("commit failed need an author and email")
	}
	author := db.NewAuthor(name, email, time.Now())
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if message == "" {
		fmt.Fprintf(os.Stderr, "commit failed, %v \r\n", "need a commit message")
	}

	tree_entries := []*db.Entry{}
	for _, index_entry := range utils.SliceFromMap(repo.Index.Entries) {
		tree_entries = append(tree_entries, db.NewEntry(index_entry.Path, index_entry.Oid, os.FileMode(index_entry.Mode)))
	}
	root := db.NewTree("")
	root.BuildTree(tree_entries)
	err = root.Traverse(func(t *db.Tree) error {
		err := repo.Database.StoreTree(t)
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
	err = repo.Database.StoreCommit(commit)
	if err != nil {
		return err
	}

	err = repo.Refs.UpdateHead(commit.Oid)
	if err != nil {
		return err
	}

	head_file, err := os.OpenFile(path.Join(repo.GitPath, "HEAD"), os.O_WRONLY|os.O_CREATE, 0777)
	defer head_file.Close()
	head_file.Write([]byte(commit.Oid))

	fmt.Println("commit:", commit.Oid, message)
	return nil
}
