package integration_test

import (
	lib "geo-git/test/integration/lib"
	"os"
	"path"
	"testing"
)

func TestIndex(t *testing.T) {
	folder := lib.GenerateRandomString()
	txt_file := lib.GenerateRandomString()
	txt_file_path := path.Join(folder, txt_file)
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	file, err := os.Create(txt_file_path)
	if err != nil {
		t.Errorf("Can't create the file on path %s", txt_file_path)
	}
	defer file.Close()
	file.Write([]byte("Test Text"))
	lib.RunGitCommand(folder, "add", ".")
	_, err = os.Stat(path.Join(folder, ".git/index"))
	if err != nil {
		t.Errorf("Index does not exist")
	}
}
