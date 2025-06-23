package commands

import (
	"fmt"
	"geo-git/lib"
	"geo-git/lib/database"
	db "geo-git/lib/database"
	"os"
	"path"
	"slices"
	"sort"
)

type Status int

const (
	WorkspaceDeleted Status = iota
	WorkspaceModified
	IndexAdded
)

type StatusTracking struct {
	Changed   []string
	Changes   map[string][]Status
	Untracked []string
	Stats     map[string]os.FileInfo
	HeadTree  map[string]*db.Entry
}

func (st *StatusTracking) Sort() {
	sort.Strings(st.Changed)
	sort.Strings(st.Untracked)
}

func statusForPath(path string, statusTracking *StatusTracking) string {
	change := statusTracking.Changes[path]
	left := ""
	right := ""

	if slices.Contains(change, WorkspaceDeleted) {
		right = "D"
	}
	if slices.Contains(change, WorkspaceModified) {
		right = "M"
	}
	if slices.Contains(change, IndexAdded) {
		left = "A"
	}
	return left + right
}

func RunStatus(repo *lib.Respository, cmd *Command) error {
	err := repo.Index.LoadForUpdate()
	if err != nil {
		return err
	}
	statusTracking := &StatusTracking{Untracked: []string{}, Changed: []string{}, Changes: map[string][]Status{}, Stats: map[string]os.FileInfo{}, HeadTree: map[string]*db.Entry{}}
	err = scanWorkspace(repo, "", statusTracking)
	if err != nil {
		return err
	}

	err = loadHeadTree(repo, "", statusTracking)
	if err != nil {
		return err
	}

	err = detectChanges(repo, "", statusTracking)
	if err != nil {
		return err
	}

	_, err = repo.Index.WriteUpdates()
	if err != nil {
		return err
	}
	statusTracking.Sort()
	printResults(statusTracking)
	return nil
}

func recordChange(statusTracking *StatusTracking, path string, status Status) {
	statusTracking.Changed = append(statusTracking.Changed, path)
	statusTracking.Changes[path] = append(statusTracking.Changes[path], status)
}

func detectChanges(repo *lib.Respository, prefix string, statusTracking *StatusTracking) error {

	//against head tree -- detectHeadTreeChanges
	for _, entry := range repo.Index.Entries {

		_, ok := statusTracking.HeadTree[entry.Path]
		if !ok {
			recordChange(statusTracking, entry.Path, IndexAdded)
		}
	}
	//against workspace -- detectWorkSpaceChanges
	for _, entry := range repo.Index.Entries {
		stat := statusTracking.Stats[entry.Path]
		if stat == nil {
			recordChange(statusTracking, entry.Path, WorkspaceDeleted)
			continue
		} else if !entry.StatMatch(stat) {
			recordChange(statusTracking, entry.Path, WorkspaceModified)
			continue
		} else if !entry.TimesMatch(stat) {
			recordChange(statusTracking, entry.Path, WorkspaceModified)
			continue
		}
		data, err := repo.Workspace.ReadFile(entry.Path)
		if err != nil {
			return err
		}
		blob := database.NewBlob(data)
		oid := blob.HashObject()
		if entry.Oid == oid {
			repo.Index.UpdateEntryStat(entry, stat)
		} else {
			recordChange(statusTracking, entry.Path, WorkspaceModified)
			continue
		}
	}
	return nil
}

func scanWorkspace(repo *lib.Respository, prefix string, statusTracking *StatusTracking) error {
	files, err := repo.Workspace.ListDirs(prefix)
	if err != nil {
		return err
	}
	for file, fileInfo := range files {
		trackable, err := trackableFile(repo, file, fileInfo)
		if err != nil {
			return err
		}
		if repo.Index.IsEntryTracked(file) {
			if fileInfo.IsDir() {
				err := scanWorkspace(repo, file, statusTracking)
				if err != nil {
					return err
				}
			} else {
				statusTracking.Stats[file] = fileInfo
			}
		} else if trackable {
			if fileInfo.IsDir() {
				file = file + string(os.PathSeparator)
			}
			statusTracking.Untracked = append(statusTracking.Untracked, file)
		}
	}
	return nil
}

func trackableFile(repo *lib.Respository, filepath string, stat os.FileInfo) (bool, error) {
	if stat == nil {
		return false, nil
	}
	if !stat.IsDir() {
		return !repo.Index.IsEntryTracked(filepath), nil
	} else {
		//depth first search
		items, err := repo.Workspace.ListDirs(filepath)
		if err != nil {
			return false, err
		}
		for item_path, item_info := range items {
			trackable, err := trackableFile(repo, item_path, item_info)
			if err != nil {
				return false, err
			}
			if trackable {
				return true, nil
			}
		}
	}
	return false, nil
}

func printResults(statusTracking *StatusTracking) {
	for path, _ := range statusTracking.Changes {
		status := statusForPath(path, statusTracking)
		fmt.Printf("%s %s\r\n", status, path)
	}
	for _, file := range slices.Compact(statusTracking.Untracked) {
		fmt.Printf("?? %s\r\n", file)
	}
}

func loadHeadTree(repo *lib.Respository, filepath string, statusTracking *StatusTracking) error {
	headOid, err := repo.Refs.ReadHead()
	if err != nil {
		return err
	}
	if headOid == "" {
		return nil
	}

	err = repo.Database.Load(headOid)
	if err != nil {
		return err
	}
	blob_commit := repo.Database.Objects[headOid]
	commit, err := db.ParseCommitFromBlob(blob_commit)
	if err != nil {
		return err
	}
	readTree(repo, commit.Tree_Oid, "", statusTracking)
	return nil

}

func readTree(repo *lib.Respository, oid string, prefix string, statusTracking *StatusTracking) error {
	if oid == "" {
		return nil
	}

	repo.Database.Load(oid)
	blob_tree := repo.Database.Objects[oid]
	tree, err := db.ParseTreeFromBlob(blob_tree)
	if err != nil {
		return err
	}

	for key, entry := range tree.Entries {
		path := path.Join(prefix, key)
		switch temp_entry := entry.(type) {
		case *db.Tree:
			err := showTree(repo, temp_entry.Oid, path)
			if err != nil {
				return err
			}
		case *db.Entry:
			statusTracking.HeadTree[path] = temp_entry
		}
	}
	return nil
}
