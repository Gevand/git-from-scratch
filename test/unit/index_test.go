package unit_test

import (
	"encoding/hex"
	"geo-git/lib"
	"os"
	"path"
	"testing"
)

func TestIndex_SingleFile(t *testing.T) {
	alice_content := []byte("Temp bytes")
	alice, err := os.Create("alice.txt")
	if err != nil {
		t.Errorf("Can't create a test file")
	}
	alice.Write(alice_content)
	alice.Close()
	defer os.Remove("alice.txt")
	tmp_path := "tmp"
	index_path := path.Join(tmp_path, "index")
	index := lib.NewIndex(index_path)
	stat, err := os.Stat("alice.txt")
	if err != nil {
		t.Errorf("Couldn't get file stat")
	}
	index.Add("alice.txt", hex.EncodeToString(alice_content), stat)
	if index.Entries["alice.txt"].Path != "alice.txt" {
		t.Errorf("Invalid path in the entries")
	}
}

func TestIndex_NextedFile(t *testing.T) {
	alice_content := []byte("Temp bytes")
	alice, err := os.Create("alice.txt")
	if err != nil {
		t.Errorf("Can't create a test file")
	}
	alice.Write(alice_content)
	alice.Close()
	defer os.Remove("alice.txt")

	bob_content := []byte("Temp bytes 2")
	bob, err := os.Create("bob.txt")
	if err != nil {
		t.Errorf("Can't create a test file")
	}
	bob.Write(bob_content)
	bob.Close()
	defer os.Remove("bob.txt")

	alice_nested_path := path.Join("alice", "alice_nexted.txt")
	os.Mkdir("alice", 0700)
	alice_nested, err := os.Create(alice_nested_path)
	if err != nil {
		t.Errorf("Can't create the nested file")
	}
	alice_nested.Write(alice_content)
	alice_nested.Close()
	defer os.Remove(alice_nested_path)

	tmp_path := "tmp"
	index_path := path.Join(tmp_path, "index")
	index := lib.NewIndex(index_path)
	alice_stat, err := os.Stat("alice.txt")
	if err != nil {
		t.Errorf("Couldn't get file stat for alice")
	}

	bob_stat, err := os.Stat("bob.txt")
	if err != nil {
		t.Errorf("Couldn't get file stat for bob")
	}

	alice_nested_stat, err := os.Stat(alice_nested_path)
	if err != nil {
		t.Errorf("Couldn't get file stat for alice_nexted")
	}

	index.Add("alice.txt", hex.EncodeToString(alice_content), alice_stat)
	index.Add("bob.txt", hex.EncodeToString(bob_content), bob_stat)
	index.Add(alice_nested_path, hex.EncodeToString(alice_content), alice_nested_stat)

	if len(index.Entries) != 2 {
		t.Errorf("Nested folders aren't properly handeled in the index, expected 2 files in index, got %d", len(index.Entries))
	}
}
