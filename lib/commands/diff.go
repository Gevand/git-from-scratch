package commands

import (
	"errors"
	"fmt"
	"geo-git/lib"
	"geo-git/lib/database"
	myers "geo-git/lib/diff"
	"geo-git/lib/index"
	repostatus "geo-git/lib/repository"
	"path/filepath"
)

type diff struct {
	path string
	oid  string
	mode uint32
	data []byte
}

const NULL_OID string = "0000000000000000000000000000000000000000"
const NULL_PATH string = "/dev/null"

func diffFromNothing(path string) (*diff, error) {
	return &diff{path: path, oid: NULL_OID, mode: 0, data: []byte("")}, nil
}
func diffFromHead(path string, statusTracking *repostatus.RepositoryStatusTracking) (*diff, error) {
	entry, ok := statusTracking.HeadTree[path]
	if !ok {
		return nil, errors.New("entry not found in head tree")
	}
	a_oid := entry.Oid
	a_mode := index.ModeForStat(uint32(entry.Mode))
	err := _repo.Database.Load(entry.Oid)
	if err != nil {
		return nil, err
	}
	blob := _repo.Database.Objects[entry.Oid]
	return &diff{path: path, oid: a_oid, mode: a_mode, data: blob.Data}, nil
}
func diffFromFile(path string, statusTracking *repostatus.RepositoryStatusTracking) (*diff, error) {
	blob_data, err := _repo.Workspace.ReadFile(path)
	if err != nil {
		return nil, err
	}
	blob := database.NewBlob(blob_data)
	b_oid := blob.HashObject()
	b_mode := index.ModeForStat(uint32(statusTracking.Stats[path].Mode()))
	return &diff{path: path, oid: b_oid, mode: b_mode, data: blob.Data}, nil
}
func diffFromIndex(path string) (*diff, error) {
	entry, err := _repo.Index.EntryForPath(path)
	if err != nil {
		return nil, err
	}

	a_oid := entry.Oid
	a_mode := index.ModeForStat(entry.Mode)
	err = _repo.Database.Load(entry.Oid)
	if err != nil {
		return nil, err
	}
	blob := _repo.Database.Objects[entry.Oid]
	return &diff{path: path, oid: a_oid, mode: a_mode, data: blob.Data}, nil
}

var _repo *lib.Respository
var _pager *lib.Pager

func RunDiff(repo *lib.Respository, cmd *Command) error {
	_repo = repo
	statusTracking := repostatus.NewStatusTracking()
	err := statusTracking.GenerateStatus(repo)
	if err != nil {
		return err
	}
	_pager = lib.NewPager()
	_pager.Initialize()
	//the order is Last in first out so they actually get executed backwards
	defer _pager.Cmd.Wait()
	defer _pager.StdIn.Close()
	if len(cmd.Args) > 0 && cmd.Args[0] == "--cached" {
		diffHeadIndex(statusTracking)
	} else {
		diffIndexWorkSpace(statusTracking)
	}

	return nil
}
func diffHeadIndex(statusTracking *repostatus.RepositoryStatusTracking) error {
	for path, state := range statusTracking.IndexChanges {
		switch state {
		case repostatus.Modified:
			a, err := diffFromHead(path, statusTracking)
			if err != nil {
				return err
			}
			b, err := diffFromIndex(path)
			if err != nil {
				return err
			}
			printDiff(*a, *b)
		}
	}
	return nil
}

func diffIndexWorkSpace(statusTracking *repostatus.RepositoryStatusTracking) error {
	for path, state := range statusTracking.WorkSpaceChanges {
		switch state {
		case repostatus.Modified:
			a, err := diffFromIndex(path)
			if err != nil {
				return err
			}
			b, err := diffFromFile(path, statusTracking)
			if err != nil {
				return err
			}
			printDiff(*a, *b)

		case repostatus.Deleted:
			a, err := diffFromIndex(path)
			if err != nil {
				return err
			}
			b, err := diffFromNothing(path)
			if err != nil {
				return nil
			}
			printDiff(*a, *b)
		}
	}
	return nil
}

func printDiff(a, b diff) {
	if a.oid == b.oid && a.mode == b.mode {
		return
	}
	a.path = filepath.Join("a", a.path)
	b.path = filepath.Join("b", b.path)
	fmt.Fprintf(_pager.StdIn, "diff --git %v %v\r\n", a.path, b.path)
	printDiffMode(a, b)
	printDiffContent(a, b)
}

func printDiffMode(a, b diff) {
	if b.mode == 0 {
		fmt.Fprintf(_pager.StdIn, "deleted file mode %d", a.mode)
	} else {
		fmt.Fprintf(_pager.StdIn, "old mode %d\n", a.mode)
		fmt.Fprintf(_pager.StdIn, "new mode %d\n", b.mode)
	}

}
func printDiffContent(a, b diff) {

	oid_range := fmt.Sprintf("index %v..%v", _repo.Database.ShortOid(a.oid), _repo.Database.ShortOid(b.oid))
	if a.mode == b.mode {
		fmt.Fprintf(_pager.StdIn, "%v %d\r\n", oid_range, a.mode)
	} else {
		fmt.Fprintf(_pager.StdIn, "%v\r\n", oid_range)
	}
	fmt.Fprintf(_pager.StdIn, "--- %v\r\n", a.path)
	fmt.Fprintf(_pager.StdIn, "+++ %v\r\n", b.path)

	myersDiff := myers.NewMyersDiff(lib.NewDiff(string(a.data), string(b.data)))
	myersDiff.DoDiff()
	hunks := myersDiff.DiffHunks()
	for _, hunk := range hunks {
		printDiffHunk(hunk)
	}
}

func printDiffHunk(hunk lib.Hunk) {
	fmt.Println(hunk.GenerateHeader())
	for _, edit := range hunk.Edits {
		line := edit.ToString()
		if line != "" {
			fmt.Println(line)
		}
	}
}
