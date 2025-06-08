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
	blob := repo.Database.Objects[oid]
	commit, err := db.ParseCommitFromBlob(blob)
	if err != nil {
		return err
	}
	fmt.Println("My Commit", commit)
	return nil
}
