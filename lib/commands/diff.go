package commands

import (
	"fmt"
	"geo-git/lib"
	"geo-git/lib/database"
	"geo-git/lib/index"
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
			err := diffFileModified(repo, statusTracking, path)
			if err != nil {
				return err
			}
		case repostatus.Deleted:
			err := diffFileDeleted(repo, statusTracking, path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func diffFileModified(repo *lib.Respository, statusTracking *repostatus.RepositoryStatusTracking, path string) error {
	entry, err := repo.Index.EntryForPath(path)
	if err != nil {
		return err
	}

	a_oid := entry.Oid
	a_mode := index.ModeForStat(entry.Mode)
	a_path := filepath.Join("a", path)

	blob_data, err := repo.Workspace.ReadFile(path)
	if err != nil {
		return err
	}
	blob := database.NewBlob(blob_data)
	b_oid := blob.HashObject()
	b_mode := index.ModeForStat(uint32(statusTracking.Stats[path].Mode()))
	b_path := filepath.Join("b", path)

	fmt.Printf("diff --git %v %v\r\n", a_path, b_path)
	if a_mode != b_mode {
		fmt.Printf("old mode %d\n", a_mode)
		fmt.Printf("new mode %d\n", b_mode)
	}

	if a_oid == b_oid {
		return nil
	}

	oid_range := fmt.Sprintf("index %v..%v", repo.Database.ShortOid(a_oid), repo.Database.ShortOid(b_oid))
	if a_mode == b_mode {
		fmt.Printf("%v %d\r\n", oid_range, a_mode)
	} else {
		fmt.Printf("%v\r\n", oid_range)
	}
	fmt.Printf("--- %v\r\n", a_path)
	fmt.Printf("+++ %v\r\n", b_path)
	return nil

}

const NULL_OID string = "0000000000000000000000000000000000000000"
const NULL_PATH string = "/dev/null"

func diffFileDeleted(repo *lib.Respository, statusTracking *repostatus.RepositoryStatusTracking, path string) error {
	entry, err := repo.Index.EntryForPath(path)
	if err != nil {
		return err
	}

	a_oid := entry.Oid
	a_mode := index.ModeForStat(entry.Mode)
	a_path := filepath.Join("a", path)

	b_oid := NULL_OID
	b_path := filepath.Join("b", path)

	fmt.Printf("diff --git %v %v\r\n", a_path, b_path)
	fmt.Printf("deleted file mode %d", a_mode)
	oid_range := fmt.Sprintf("index %v..%v", repo.Database.ShortOid(a_oid), repo.Database.ShortOid(b_oid))
	fmt.Printf("%v %d\r\n", oid_range, a_mode)
	fmt.Printf("--- %v\r\n", a_path)
	fmt.Printf("+++ %v\r\n", NULL_PATH)
	return nil
}
