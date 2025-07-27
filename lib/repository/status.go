package repository

import (
	"geo-git/lib"
	db "geo-git/lib/database"
	ind "geo-git/lib/index"
	"os"
	"path"
	"sort"
)

type Status int

const (
	Deleted Status = iota
	Modified
	Added
)

type RepositoryStatusTracking struct {
	Changed          []string
	IndexChanges     map[string]Status
	WorkSpaceChanges map[string]Status
	Untracked        []string
	Stats            map[string]os.FileInfo
	HeadTree         map[string]*db.Entry
}

func NewStatusTracking() *RepositoryStatusTracking {
	return &RepositoryStatusTracking{Untracked: []string{}, Changed: []string{}, IndexChanges: map[string]Status{}, WorkSpaceChanges: map[string]Status{}, Stats: map[string]os.FileInfo{}, HeadTree: map[string]*db.Entry{}}

}

var statusRepo *lib.Respository

func (statusTracking *RepositoryStatusTracking) GenerateStatus(repo *lib.Respository) error {
	statusRepo = repo
	err := statusRepo.Index.LoadForUpdate()
	if err != nil {
		return err
	}
	err = scanWorkspace("", statusTracking)
	if err != nil {
		return err
	}

	err = loadHeadTree(statusTracking)
	if err != nil {
		return err
	}

	err = checkIndexEntries(statusRepo, statusTracking)
	if err != nil {
		return err
	}

	_, err = statusRepo.Index.WriteUpdates()
	if err != nil {
		return err
	}
	statusTracking.Sort()
	return nil
}

func (st *RepositoryStatusTracking) Sort() {
	sort.Strings(st.Changed)
	sort.Strings(st.Untracked)
}

func recordChange(statusTracking *RepositoryStatusTracking, path string, changeMap map[string]Status, status Status) {
	statusTracking.Changed = append(statusTracking.Changed, path)
	changeMap[path] = status
}

func checkIndexEntries(repo *lib.Respository, statusTracking *RepositoryStatusTracking) error {

	for _, entry := range repo.Index.Entries {
		err := checkIndexAgainstWorkSpace(entry, statusTracking)
		if err != nil {
			return err
		}

		err = checkIndexAgainstHeadTree(entry, statusTracking)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkIndexAgainstWorkSpace(entry *ind.IndexEntry, statusTracking *RepositoryStatusTracking) error {
	stat := statusTracking.Stats[entry.Path]
	if stat == nil {
		recordChange(statusTracking, entry.Path, statusTracking.WorkSpaceChanges, Deleted)
		return nil
	} else if !entry.StatMatch(stat) {
		recordChange(statusTracking, entry.Path, statusTracking.WorkSpaceChanges, Modified)
		return nil

	} else if !entry.TimesMatch(stat) {
		recordChange(statusTracking, entry.Path, statusTracking.WorkSpaceChanges, Modified)
		return nil
	}
	data, err := statusRepo.Workspace.ReadFile(entry.Path)
	if err != nil {
		return err
	}
	blob := db.NewBlob(data)
	oid := blob.HashObject()
	if entry.Oid == oid {
		statusRepo.Index.UpdateEntryStat(entry, stat)
	} else {
		recordChange(statusTracking, entry.Path, statusTracking.WorkSpaceChanges, Modified)
		return nil
	}

	return nil
}

func checkIndexAgainstHeadTree(entry *ind.IndexEntry, statusTracking *RepositoryStatusTracking) error {
	found_item, ok := statusTracking.HeadTree[entry.Path]
	if ok && found_item != nil {
		if found_item.Mode != os.FileMode(entry.Mode) || found_item.Oid != entry.Oid {
			recordChange(statusTracking, entry.Path, statusTracking.IndexChanges, Modified)
		}
	} else if !ok {
		recordChange(statusTracking, entry.Path, statusTracking.IndexChanges, Added)
	}

	for path, _ := range statusTracking.HeadTree {
		if !statusRepo.Index.IsEntryTracked(path) {
			recordChange(statusTracking, path, statusTracking.IndexChanges, Deleted)
		}
	}

	return nil
}

func scanWorkspace(prefix string, statusTracking *RepositoryStatusTracking) error {
	files, err := statusRepo.Workspace.ListDirs(prefix)
	if err != nil {
		return err
	}
	for file, fileInfo := range files {
		trackable, err := trackableFile(file, fileInfo)
		if err != nil {
			return err
		}
		if statusRepo.Index.IsEntryTracked(file) {
			if fileInfo.IsDir() {
				err := scanWorkspace(file, statusTracking)
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

func loadHeadTree(statusTracking *RepositoryStatusTracking) error {
	headOid, err := statusRepo.Refs.ReadHead()
	if err != nil {
		return err
	}
	if headOid == "" {
		return nil
	}

	err = statusRepo.Database.Load(headOid)
	if err != nil {
		return err
	}
	blob_commit := statusRepo.Database.Objects[headOid]
	commit, err := db.ParseCommitFromBlob(blob_commit)
	if err != nil {
		return err
	}
	readTree(commit.Tree_Oid, "", statusTracking)
	return nil

}

func readTree(oid string, prefix string, statusTracking *RepositoryStatusTracking) error {
	if oid == "" {
		return nil
	}

	statusRepo.Database.Load(oid)
	blob_tree := statusRepo.Database.Objects[oid]
	tree, err := db.ParseTreeFromBlob(blob_tree)
	if err != nil {
		return err
	}

	for key, entry := range tree.Entries {
		path := path.Join(prefix, key)
		switch temp_entry := entry.(type) {
		case *db.Tree:
			err := readTree(temp_entry.Oid, path, statusTracking)
			if err != nil {
				return err
			}
		case *db.Entry:
			statusTracking.HeadTree[path] = temp_entry
		}
	}
	return nil
}

func trackableFile(filepath string, stat os.FileInfo) (bool, error) {
	if stat == nil {
		return false, nil
	}
	if !stat.IsDir() {
		return !statusRepo.Index.IsEntryTracked(filepath), nil
	} else {
		//depth first search
		items, err := statusRepo.Workspace.ListDirs(filepath)
		if err != nil {
			return false, err
		}
		for item_path, item_info := range items {
			trackable, err := trackableFile(item_path, item_info)
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
