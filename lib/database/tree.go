package database

import (
	"fmt"
	"geo-git/lib/utils"
	"path/filepath"
	"sort"
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
			return_value = fmt.Sprintf("%s\000%s", temp_string, oid_as_hexstring)
		case *Tree:
			temp_string := fmt.Sprintf("%v %v\000", fmt.Sprintf("%06o", DIRECTORY_MODE), entry.Name)
			oid_as_hexstring := string(utils.PackHexaDecimal(entry.Oid))
			return_value += temp_string
			return_value += oid_as_hexstring
		}
	}
	return return_value
}
