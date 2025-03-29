package lib

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	ind "geo-git/lib/index"
	"geo-git/lib/utils"
	"hash"
	"os"
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
	changed  bool
}

func NewIndex(path string) *Index {
	return &Index{
		path:     path,
		Entries:  map[string]*ind.IndexEntry{},
		lockfile: NewLockFile(path),
		keys:     []string{},
		changed:  false,
	}
}

func (i *Index) Add(path, oid string, stat os.FileInfo) error {
	index_entry, err := ind.NewEntry(stat, path, oid)
	if err != nil {
		return err
	}
	i.StoreEntry(index_entry)
	i.changed = true
	return err
}

func (i *Index) StoreEntry(entry *ind.IndexEntry) {
	i.Entries[entry.Path] = entry
	i.keys = append(i.keys, entry.Path)
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

func (i *Index) LoadForUpdate() (bool, error) {
	err := i.lockfile.HoldForUpdate()
	if err == nil {
		err := i.Load()
		return err == nil, err
	}
	return false, nil
}

func (i *Index) Load() error {
	index_file, err := os.Open(i.path)
	defer index_file.Close()
	if err != nil {
		return err
	}

	reader := ind.NewChecksum(*index_file)
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
		for {
			entry_block := make([]byte, ind.ENTRY_MIN_SIZE)
			_, err := reader.File.Read(entry_block)
			if err != nil {
				return err
			}
			entry_bytes = append(entry_bytes, entry_block...)
			if entry_bytes[len(entry_bytes)-1] == 0 {
				break
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
