package commands

import (
	"fmt"
	"geo-git/lib"
	db "geo-git/lib/database"
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
	repo.Database.Load(commit.Tree_Oid)
	blob_tree := repo.Database.Objects[commit.Tree_Oid]
	tree, err := db.ParseTreeFromBlob(blob_tree)
	if err != nil {
		return err
	}
	fmt.Println("My Tree", tree)
	//todo: parse the tree
	return nil
}
