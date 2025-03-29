package lib

import (
	"crypto/sha1"
	ind "geo-git/lib/index"
	"geo-git/lib/utils"
	"hash"
	"os"
	"sort"
)

type Index struct {
	path     string
	Entries  map[string]*ind.IndexEntry
	lockfile *LockFile
	digest   hash.Hash
	keys     []string
}

func NewIndex(path string) *Index {
	return &Index{
		path:     path,
		Entries:  map[string]*ind.IndexEntry{},
		lockfile: NewLockFile(path),
		keys:     []string{},
	}
}

func (i *Index) Add(path, oid string, stat os.FileInfo) error {
	index_entry, err := ind.NewEntry(stat, path, oid)
	if err != nil {
		return err
	}
	i.Entries[path] = index_entry
	i.keys = append(i.keys, path)
	return err
}

func (i *Index) WriteUpdates() (bool, error) {
	err := i.lockfile.HoldForUpdate()
	if err != nil {
		return false, err
	}
	i.BegindWrite()
	header := []byte{}
	header = append(header, []byte("DIRC")...)

	version_number_as_bytes, err := utils.Int32ToBigEndianBytes(2)
	if err != nil {
		return false, err
	}
	header = append(header, version_number_as_bytes...)

	entries_length_as_bytes, err := utils.Int32ToBigEndianBytes(uint32(len(i.Entries)))
	if err != nil {
		return false, err
	}
	header = append(header, entries_length_as_bytes...)
	i.Write(header)
	//keys = path, and need to be sorted
	sort.Strings(i.keys)
	for _, key := range i.keys {
		entry := i.Entries[key]
		str, err := entry.ToString()
		if err != nil {
			return false, err
		}
		i.Write([]byte(str))
	}
	err = i.FinishWrite()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (i *Index) BegindWrite() {
	i.digest = sha1.New()
}

func (i *Index) Write(data []byte) error {
	err := i.lockfile.Write(data)
	if err != nil {
		return err
	}
	_, err = i.digest.Write(data)
	return err
}

func (i *Index) FinishWrite() error {
	digest := i.digest.Sum(nil)
	err := i.lockfile.Write(digest)
	if err != nil {
		return err
	}
	return i.lockfile.Commit()
}
