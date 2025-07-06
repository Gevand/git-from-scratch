package database

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"geo-git/lib/utils"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Tree struct {
	Entries map[string]interface{}
	Oid     string
	Name    string
}

func NewTree(name string) *Tree {
	root := &Tree{Entries: map[string]interface{}{}, Name: name}
	return root
}

func (t *Tree) BuildTree(entries []*Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	for _, entry := range entries {
		parents := entry.ParentDirectories()
		if len(parents) > 0 {
			entry.Name = filepath.Base(entry.Name)
		}
		t.AddEntry(parents, entry)
	}
}

func (t *Tree) AddEntry(parents []string, entry *Entry) {
	if len(parents) == 0 {
		t.Entries[entry.Name] = entry
	} else {
		tree, exists := t.Entries[parents[0]]
		if !exists {
			tree = NewTree(parents[0])
		}

		tree.(*Tree).AddEntry(parents[1:], entry)
		t.Entries[parents[0]] = tree
	}
}

func (t *Tree) Traverse(execute func(*Tree) error) error {
	for _, interface_entry := range t.Entries {
		switch entry := interface_entry.(type) {
		case *Tree:
			entry.Traverse(execute)
		}
	}
	err := execute(t)
	return err
}

func (t *Tree) ToString() string {
	return_value := ""
	for _, interface_entry := range t.Entries {
		switch entry := interface_entry.(type) {
		case *Entry:
			temp_string := fmt.Sprintf("%v %v", fmt.Sprintf("%06o", entry.Mode), entry.Name)
			oid_as_hexstring := string(utils.PackHexaDecimal(entry.Oid))
			return_value += fmt.Sprintf("%s\000%s", temp_string, oid_as_hexstring)
		case *Tree:
			temp_string := fmt.Sprintf("%v %v", fmt.Sprintf("%06o", DIRECTORY_MODE), entry.Name)
			oid_as_hexstring := string(utils.PackHexaDecimal(entry.Oid))
			return_value += fmt.Sprintf("%s\000%s", temp_string, oid_as_hexstring)
		}
	}
	return strings.TrimSpace(return_value)
}

func ParseTreeFromBlob(blob *Blob) (*Tree, error) {
	treeToReturn := &Tree{Entries: map[string]interface{}{}}
	blob_data := string(blob.Data)
	//first line of the blob is tree space length\000, get rid of that
	blob_prefix := (strings.Split(blob_data, "\000")[0]) + "\000"
	blob_data = strings.Replace(blob_data, blob_prefix, "", 1)
	//starts the parsing process
	reader := bufio.NewReader(strings.NewReader(blob_data))
	for {
		mode_bytes, err := reader.ReadBytes(' ')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		mode_bytes = mode_bytes[:len(mode_bytes)-1] // git rid of the ' ' byte at the end

		name_bytes, err := reader.ReadBytes('\000')
		if err != nil {
			return nil, err
		}
		name_bytes = name_bytes[:len(name_bytes)-1] //get rid of the 0 byte at the end

		oid_bytes, err := reader.Peek(20)
		if err != nil {
			return nil, err
		}

		_, err = reader.Discard(20)
		if err != nil {
			return nil, err
		}

		var last_entry interface{}
		mode := string(mode_bytes)
		name := string(name_bytes)
		oid := hex.EncodeToString(oid_bytes)
		if mode == fmt.Sprintf("%06o", DIRECTORY_MODE) {
			last_entry = &Tree{Name: name, Entries: map[string]interface{}{}, Oid: oid}
		} else {
			mode, err := strconv.ParseUint(mode, 8, 32)
			if err != nil {
				return nil, err
			}
			last_entry = &Entry{Name: name, Mode: os.FileMode(uint32(mode)), Oid: oid}
		}
		treeToReturn.Entries[name] = last_entry
	}

	return treeToReturn, nil
}
