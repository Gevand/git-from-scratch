package database

import (
	"encoding/hex"
	"errors"
	"fmt"
	"geo-git/lib/utils"
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
	entry_parts := strings.Split(string(blob.Data), "\000")

	entry_name := ""
	entry_mode := ""
	var last_entry interface{}
	for index, entry_part := range entry_parts {

		if index == 0 {
			//first entry is always "%v %v"
			temp_split := strings.Split(entry_part, " ")
			if len(temp_split) != 2 {
				return nil, errors.New("blob is not formatted as a proper tree")
			}
			entry_mode = temp_split[0]
			entry_name = temp_split[1]
			if entry_mode == fmt.Sprintf("%06o", DIRECTORY_MODE) {
				last_entry = &Tree{Name: entry_name, Entries: map[string]interface{}{}}
			} else {

				mode, err := strconv.ParseUint(entry_mode, 8, 32)
				if err != nil {
					return nil, err
				}
				last_entry = &Entry{Name: entry_name, Mode: os.FileMode(uint32(mode))}
			}
			treeToReturn.Entries[entry_name] = last_entry
		} else if index == len(entry_parts)-1 {
			//last entry is always "%s"
			previous_oid := hex.EncodeToString([]byte(entry_part))
			switch entry := last_entry.(type) {
			case *Entry:
				entry.Oid = previous_oid
			case *Tree:
				entry.Oid = previous_oid
			}
		} else {
			//everything else is "%s%v %v"
			previous_oid := hex.EncodeToString([]byte(entry_part)[0:20])
			switch entry := last_entry.(type) {
			case *Entry:
				entry.Oid = previous_oid
			case *Tree:
				entry.Oid = previous_oid
			}
			temp_split := strings.Split(entry_part[20:], " ")
			if len(temp_split) != 2 {
				return nil, errors.New("blob is not formatted as a proper tree")
			}
			entry_mode = temp_split[0]
			entry_name = temp_split[1]
			if entry_mode == fmt.Sprintf("%06o", DIRECTORY_MODE) {
				last_entry = &Tree{Name: entry_name, Entries: map[string]interface{}{}}
			} else {

				mode, err := strconv.ParseUint(entry_mode, 8, 32)
				if err != nil {
					return nil, err
				}
				last_entry = &Entry{Name: entry_name, Mode: os.FileMode(uint32(mode))}
			}
			treeToReturn.Entries[entry_name] = last_entry
		}
	}
	return treeToReturn, nil
}
