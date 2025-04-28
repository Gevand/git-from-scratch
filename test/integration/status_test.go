package integration

import (
	lib "geo-git/test/integration/lib"
	"os"
	"path"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
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
	status_output := lib.RunGitCommandWithOutput(folder, "status")
	if !strings.Contains(status_output, txt_file) {
		t.Errorf("Status command didn't return the expected output")
	}
}
