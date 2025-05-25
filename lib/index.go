package lib

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	ind "geo-git/lib/index"
	"geo-git/lib/utils"
	"hash"
	"os"
	"slices"
	"sort"
)

const (
	SIGNATURE   = "DIRC"
	VERSION     = 2
	HEADER_SIZE = 12
)

type Index struct {
	path     string
	Entries  map[string]*ind.IndexEntry
	lockfile *LockFile
	digest   hash.Hash
	keys     []string
	parents  map[string][]string
	changed  bool
}

func NewIndex(path string) *Index {
	return &Index{
		path:     path,
		Entries:  map[string]*ind.IndexEntry{},
		lockfile: NewLockFile(path),
		keys:     []string{},
		parents:  map[string][]string{},
		changed:  false,
	}
}

func (i *Index) Add(path, oid string, stat os.FileInfo) error {
	index_entry, err := ind.NewEntry(stat, path, oid)
	if err != nil {
		return err
	}
	i.DiscardConflicts(index_entry)
	i.StoreEntry(index_entry)
	i.changed = true
	return err
}

func (i *Index) DiscardConflicts(entry *ind.IndexEntry) {
	for _, dirname := range entry.ParentDirectories() {
		i.removeEntry(dirname)
	}
	i.removeChildren(entry.Path)
}

func (i *Index) removeChildren(path string) {
	paths, ok := i.parents[path]
	if !ok {
		return
	}
	for _, p := range paths {
		i.removeEntry(p)
	}
}

func (i *Index) removeEntry(path string) {
	entry, ok := i.Entries[path]
	if !ok {
		return
	}
	i.keys = slices.DeleteFunc(i.keys, func(s string) bool { return s == entry.Path })
	delete(i.Entries, entry.Path)
	for _, dirname := range entry.ParentDirectories() {
		i.parents[dirname] = slices.DeleteFunc(i.parents[dirname], func(s string) bool { return s == entry.Path })
		if len(i.parents[dirname]) == 0 {
			delete(i.parents, dirname)
		}
	}
}

func (i *Index) StoreEntry(entry *ind.IndexEntry) {
	i.Entries[entry.Path] = entry
	i.keys = append(i.keys, entry.Path)

	for _, dirname := range entry.ParentDirectories() {
		i.parents[dirname] = append(i.parents[dirname], entry.Path)
	}
}

func (i *Index) IsEntryTracked(path string) bool {
	_, tracked_file := i.Entries[path]
	_, tracked_parent := i.parents[path]
	return tracked_file || tracked_parent
}

func (i *Index) WriteUpdates() (bool, error) {
	err := i.lockfile.HoldForUpdate()
	if err != nil {
		return false, err
	}
	i.BegindWrite()
	header := []byte{}
	header = append(header, []byte(SIGNATURE)...)

	version_number_as_bytes, err := utils.Int32ToBigEndianBytes(VERSION)
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

func (i *Index) LoadForUpdate() error {
	err := i.lockfile.HoldForUpdate()
	if err == nil {
		err := i.Load()
		return err
	}
	return nil
}

func (i *Index) Load() error {

	index_file, err := os.Open(i.path)
	defer index_file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	reader := ind.NewChecksum(index_file)
	count, err := i.ReadHeader(reader)
	if err != nil {
		return err
	}

	err = i.ReadEntries(reader, count)
	if err != nil {
		return err
	}

	return reader.Verify()

}

func (i *Index) ReadHeader(reader *ind.Checksum) (int, error) {
	data, err := reader.Read(HEADER_SIZE)
	if err != nil {
		return 0, err
	}
	signature := string(data[:4])
	if signature != SIGNATURE {
		return 0, errors.New("Invalid signature, expcted " + SIGNATURE + " got signature")
	}
	version := binary.BigEndian.Uint32(data[4:8])
	if version != VERSION {
		return 0, errors.New("Invalid version, expected" + string(rune(VERSION)) + " got " + string(rune(version)))
	}

	count := binary.BigEndian.Uint32(data[8:])

	return int(count), nil
}

func (i *Index) ReadEntries(reader *ind.Checksum, count int) error {
	for c := 0; c < count; c++ {
		entry_bytes := []byte{}
		entry_block, err := reader.Read(ind.ENTRY_MIN_SIZE)
		if err != nil {
			return err
		}
		for {

			entry_bytes = append(entry_bytes, entry_block...)
			if entry_bytes[len(entry_bytes)-1] == 0 {
				break
			}
			entry_block, err = reader.Read(ind.ENTRY_BLOCK)
			if err != nil {
				return err
			}
		}
		entry, err := ind.ParseEntry(entry_bytes)
		if err != nil {
			return err
		}
		i.StoreEntry(entry)
	}
	return nil
}

func (i *Index) UpdateEntryStat(entry *ind.IndexEntry, stat os.FileInfo) {
	entry.UpdateStat(stat)
	i.changed = true
}
