package lib

import (
	ind "geo-git/lib/index"
	"geo-git/lib/utils"
	"os"
)

type Index struct {
	path     string
	Entries  map[string]*ind.IndexEntry
	lockfile *LockFile
}

func NewIndex(path string) *Index {
	return &Index{
		path:     path,
		Entries:  map[string]*ind.IndexEntry{},
		lockfile: NewLockFile(path),
	}
}

func (i *Index) Add(path, oid string, stat os.FileInfo) error {
	index_entry, err := ind.NewEntry(stat, path, oid)
	if err != nil {
		return err
	}
	i.Entries[path] = index_entry
	return err
}

func (i *Index) WriteUpdates() (bool, error) {
	err := i.lockfile.HoldForUpdate()
	if err != nil {
		return false, err
	}
	i.BegindWrite()
	defer i.FinishWrite()
	header := []byte{}
	header = append(header, []byte("DIRC")...)
	header = append(header, byte(2))
	entries_length_as_bytes, err := utils.NumberToBigEndianBytes(uint32(len(i.Entries)))
	if err != nil {
		return false, err
	}
	header = append(header, entries_length_as_bytes...)
	i.Write(header)
	for _, entry := range i.Entries {
		i.Write([]byte(entry.ToString()))
	}
	return true, nil
}

func (i *Index) BegindWrite() {}
func (i *Index) Write(data []byte) error {
	return nil
}
func (i *Index) FinishWrite() {}
