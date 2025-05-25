package commands

import (
	"fmt"
	"geo-git/lib"
	"geo-git/lib/database"
	"os"
	"slices"
	"sort"
)

func RunStatus(repo *lib.Respository, cmd *Command) error {
	err := repo.Index.LoadForUpdate()
	if err != nil {
		return err
	}
	untracked := []string{}
	stats := map[string]os.FileInfo{}
	err = scanWorkspace(repo, "", &untracked, &stats)
	if err != nil {
		return err
	}

	changed := []string{}
	err = detectWorkspaceChanges(repo, "", &changed, &stats)
	if err != nil {
		return err
	}

	_, err = repo.Index.WriteUpdates()
	if err != nil {
		return err
	}

	sort.Strings(changed)
	for _, file := range slices.Compact(changed) {
		fmt.Printf("M %s\r\n", file)
	}

	sort.Strings(untracked)
	for _, file := range slices.Compact(untracked) {
		fmt.Printf("?? %s\r\n", file)
	}
	return nil
}

func detectWorkspaceChanges(repo *lib.Respository, prefix string, changed *[]string, stats *map[string]os.FileInfo) error {
	for _, entry := range repo.Index.Entries {
		stat := (*stats)[entry.Path]
		if !entry.StatMatch(stat) {
			*changed = append(*changed, entry.Path)
			continue
		} else if !entry.TimesMatch(stat) {
			*changed = append(*changed, entry.Path)
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
			*changed = append(*changed, entry.Path)
			continue
		}
	}
	return nil
}

func scanWorkspace(repo *lib.Respository, prefix string, untracked *[]string, stats *map[string]os.FileInfo) error {
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
				err := scanWorkspace(repo, file, untracked, stats)
				if err != nil {
					return err
				}
			} else {
				(*stats)[file] = fileInfo
			}
		} else if trackable {
			if fileInfo.IsDir() {
				file = file + string(os.PathSeparator)
			}
			*untracked = append(*untracked, file)
		}
	}
	return nil
}

func trackableFile(repo *lib.Respository, file_path string, stat os.FileInfo) (bool, error) {
	if stat == nil {
		return false, nil
	}
	if !stat.IsDir() {
		return !repo.Index.IsEntryTracked(file_path), nil
	} else {
		//depth first search
		items, err := repo.Workspace.ListDirs(file_path)
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
