package commands

import (
	"fmt"
	"geo-git/lib"
	"geo-git/lib/database"
	repostatus "geo-git/lib/repository"
	"path/filepath"
)

func RunDiff(repo *lib.Respository, cmd *Command) error {
	statusTracking := repostatus.NewStatusTracking()
	err := statusTracking.GenerateStatus(repo)
	if err != nil {
		return err
	}
	for path, state := range statusTracking.WorkSpaceChanges {
		switch state {
		case repostatus.Modified:
			err := diffFile(repo, statusTracking, path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func diffFile(repo *lib.Respository, statusTracking *repostatus.RepositoryStatusTracking, path string) error {
	entry, err := repo.Index.EntryForPath(path)
	if err != nil {
		return err
	}

	a_oid := entry.Oid
	a_mode := entry.Mode
	a_path := filepath.Join("a", path)

	blob_data, err := repo.Workspace.ReadFile(path)
	if err != nil {
		return err
	}
	blob := database.NewBlob(blob_data)
	b_oid := blob.HashObject()
	b_path := filepath.Join("b", path)

	fmt.Printf("diff --git %v %v\r\n", a_path, b_path)
	fmt.Printf("index %v..%v %06o\r\n", repo.Database.ShortOid(a_oid), repo.Database.ShortOid(b_oid), a_mode)
	fmt.Printf("--- %v\r\n", a_path)
	fmt.Printf("+++ %v\r\n", b_path)
	return nil

}
