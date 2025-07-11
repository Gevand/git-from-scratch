package commands

import (
	"fmt"
	"geo-git/lib"
	"geo-git/lib/database"
	db "geo-git/lib/database"
	ind "geo-git/lib/index"
	"geo-git/lib/utils"
	"os"
	"path"
	"slices"
	"sort"
)

type Status int

const (
	Deleted Status = iota
	Modified
	Added
)

var ColorMap = map[string]string{
	"green": "\033[32m",
	"red":   "\033[31m",
	"reset": "\033[0m",
}

var ShortStatusMap = map[Status]string{
	Deleted:  "D",
	Modified: "M",
	Added:    "A",
}
var LongStatusMap = map[Status]string{
	Deleted:  "deleted",
	Modified: "modified",
	Added:    "new file",
}

type StatusTracking struct {
	Changed          []string
	IndexChanges     map[string]Status
	WorkSpaceChanges map[string]Status
	Untracked        []string
	Stats            map[string]os.FileInfo
	HeadTree         map[string]*db.Entry
}

func (st *StatusTracking) Sort() {
	sort.Strings(st.Changed)
	sort.Strings(st.Untracked)
}

var statusRepo *lib.Respository

func RunStatus(repo *lib.Respository, cmd *Command) error {
	statusRepo = repo
	err := statusRepo.Index.LoadForUpdate()
	if err != nil {
		return err
	}
	statusTracking := &StatusTracking{Untracked: []string{}, Changed: []string{}, IndexChanges: map[string]Status{}, WorkSpaceChanges: map[string]Status{}, Stats: map[string]os.FileInfo{}, HeadTree: map[string]*db.Entry{}}
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

	if utils.Contains(cmd.Args, "--porcelain") {
		printResultsPorcelain(statusTracking)
	} else {
		printResultsLongs(statusTracking)
	}
	return nil
}

func recordChange(statusTracking *StatusTracking, path string, changeMap map[string]Status, status Status) {
	statusTracking.Changed = append(statusTracking.Changed, path)
	changeMap[path] = status
}

func checkIndexEntries(repo *lib.Respository, statusTracking *StatusTracking) error {

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

func checkIndexAgainstWorkSpace(entry *ind.IndexEntry, statusTracking *StatusTracking) error {
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
	blob := database.NewBlob(data)
	oid := blob.HashObject()
	if entry.Oid == oid {
		statusRepo.Index.UpdateEntryStat(entry, stat)
	} else {
		recordChange(statusTracking, entry.Path, statusTracking.WorkSpaceChanges, Modified)
		return nil
	}

	return nil
}

func checkIndexAgainstHeadTree(entry *ind.IndexEntry, statusTracking *StatusTracking) error {
	found_item, ok := statusTracking.HeadTree[entry.Path]
	if ok && found_item != nil {
		if found_item.Mode != os.FileMode(entry.Mode) || found_item.Oid != entry.Oid {
			recordChange(statusTracking, entry.Path, statusTracking.IndexChanges, Modified)
		}
	} else if !ok {
		fmt.Println("DEBUG", entry.Path, " added because not found in head tree")
		recordChange(statusTracking, entry.Path, statusTracking.IndexChanges, Added)
	}

	for path, _ := range statusTracking.HeadTree {
		if !statusRepo.Index.IsEntryTracked(path) {
			recordChange(statusTracking, path, statusTracking.IndexChanges, Deleted)
		}
	}

	return nil
}

func scanWorkspace(prefix string, statusTracking *StatusTracking) error {
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

func loadHeadTree(statusTracking *StatusTracking) error {
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

func readTree(oid string, prefix string, statusTracking *StatusTracking) error {
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

func printResultsLongs(statusTracking *StatusTracking) {
	printChanges("Changes to be committed", statusTracking.IndexChanges, "green")
	printChanges("Changes not staged for commit", statusTracking.WorkSpaceChanges, "red")
	printUntrackedChanges("Untracked files", statusTracking.Untracked, "red")
	printCommitStatus(statusTracking)
}

func printCommitStatus(statusTracking *StatusTracking) {
	if len(statusTracking.IndexChanges) > 0 {
		return
	}

	if len(statusTracking.WorkSpaceChanges) > 0 {
		fmt.Println("no changes added to commit")
	} else if len(statusTracking.Untracked) > 0 {
		fmt.Println("nothing added to commit but untracked files present")
	} else {
		fmt.Println("nothing to commit, working tree clean")
	}
}

func printChanges(message string, changes map[string]Status, color string) {
	if len(changes) == 0 {
		return
	}
	fmt.Println(message)
	fmt.Println("")
	color_code, ok := ColorMap[color]
	reset := ColorMap["reset"]
	for path, status := range changes {
		if ok {
			fmt.Printf("%s%s %s\n%s", color_code, LongStatusMap[status], path, reset)
		} else {
			fmt.Printf("%s %s\n", LongStatusMap[status], path)
		}
	}
	fmt.Println("")
}
func printUntrackedChanges(message string, changes []string, color string) {
	if len(changes) == 0 {
		return
	}
	fmt.Println(message)
	fmt.Println("")
	color_code, ok := ColorMap[color]
	reset := ColorMap["reset"]
	for _, path := range changes {
		if ok {
			fmt.Printf("%s%s\n%s", color_code, path, reset)
		} else {
			fmt.Printf("%s\n", path)
		}
	}
	fmt.Println("")
}
func printResultsPorcelain(statusTracking *StatusTracking) {
	for path, _ := range statusTracking.IndexChanges {
		status := statusForPath(path, statusTracking)
		fmt.Printf("%s %s\r\n", status, path)
	}
	for _, file := range slices.Compact(statusTracking.Untracked) {
		fmt.Printf("?? %s\r\n", file)
	}
}

func statusForPath(path string, statusTracking *StatusTracking) string {
	var left string
	var right string
	left_status, ok := statusTracking.IndexChanges[path]
	if !ok {
		left = ""
	} else {
		left = ShortStatusMap[left_status]
	}

	right_status, ok := statusTracking.WorkSpaceChanges[path]
	if !ok {
		right = ""
	} else {
		right = ShortStatusMap[right_status]
	}

	return left + right
}
