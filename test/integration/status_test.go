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
	status_output := lib.RunGitCommandWithOutput(folder, "status", "--porcelain")
	if !strings.Contains(status_output, "?? "+txt_file) {
		t.Errorf("Status command didn't return the expected output")
	}
}

func TestStatus_IgnoreIndexFiles(t *testing.T) {
	folder := lib.GenerateRandomString()
	txt_file_tracked := lib.GenerateRandomString()
	txt_file_tracked_path := path.Join(folder, txt_file_tracked)
	txt_file_untracked := lib.GenerateRandomString()
	txt_file_untracked_path := path.Join(folder, txt_file_untracked)
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	file_tracked, err := os.Create(txt_file_tracked_path)
	if err != nil {
		t.Errorf("Can't create the file on path %s", txt_file_tracked_path)
	}
	defer file_tracked.Close()
	file_tracked.Write([]byte("Tracked"))
	file_untracked, err := os.Create(txt_file_untracked_path)
	if err != nil {
		t.Errorf("Can't create the file on path %s", txt_file_untracked_path)
	}
	defer file_untracked.Close()
	file_untracked.Write([]byte("Untracked"))
	lib.RunGitCommand(folder, "add", txt_file_tracked)
	status_output_porcelain := lib.RunGitCommandWithOutput(folder, "status", "--porcelain")
	if !strings.Contains(status_output_porcelain, "?? "+txt_file_untracked) {
		t.Errorf("Status command --porcelain didn't return the expected output: %s should be untracked", txt_file_tracked)
	}
	if strings.Contains(status_output_porcelain, "?? "+txt_file_tracked) {
		t.Errorf("Status command --porcelain didn't return the expected output: %s shouldn't be untracked", txt_file_tracked)
	}
}

func TestStatus_WorkSpaceChange(t *testing.T) {
	folder := lib.GenerateRandomString()
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	one_txt, err := os.Create(path.Join(folder, "1.txt"))
	if err != nil {
		t.Errorf("Can't create the file on path %s", "1.txt")
		return
	}
	defer one_txt.Close()
	one_txt.WriteString("1")

	lib.RunCustomCommand(folder, "mkdir", "a")
	two_txt, err := os.Create(path.Join(folder, "a/2.txt"))
	if err != nil {
		t.Errorf("Can't create the file on path %s", "a/2.txt")
		return
	}
	defer two_txt.Close()

	two_txt.WriteString("2")
	lib.RunGitCommand(folder, "add", ".")
	lib.RunGitCommand(folder, "commit", "some message")

	//modifying the second file
	two_txt.WriteString("Modified")
	status_output := lib.RunGitCommandWithOutput(folder, "status", "")
	if !strings.Contains(status_output, "modified "+"a/2.txt") {
		t.Errorf("Status command didn't return the expected output: %s should be modified, got %s", "a/2.txt", status_output)
	}
}

func TestStatus_WorkSpaceChange_FileIsExecutable(t *testing.T) {
	folder := lib.GenerateRandomString()
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	one_txt, err := os.Create(path.Join(folder, "1.txt"))
	if err != nil {
		t.Errorf("Can't create the file on path %s", "1.txt")
		return
	}
	defer one_txt.Close()
	one_txt.WriteString("1")

	lib.RunGitCommand(folder, "add", ".")
	lib.RunGitCommand(folder, "commit", "some message")

	//making one_txt file executable
	err = os.Chmod(path.Join(folder, "1.txt"), 0755)
	if err != nil {
		t.Errorf("Can't change the file's mode")
		return
	}

	status_output := lib.RunGitCommandWithOutput(folder, "status", "")
	if !strings.Contains(status_output, "modified "+"1.txt") {
		t.Errorf("Status command didn't return the expected output: %s should be modified, got %s", "1.txt", status_output)
	}
}

func TestStatus_WorkSpaceChange_DoubleSave(t *testing.T) {
	folder := lib.GenerateRandomString()
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	one_txt, err := os.Create(path.Join(folder, "1.txt"))
	if err != nil {
		t.Errorf("Can't create the file on path %s", "1.txt")
		return
	}
	one_txt.WriteString("1")

	lib.RunGitCommand(folder, "add", ".")
	lib.RunGitCommand(folder, "commit", "some message")

	one_txt.Close()
	//write the same content again
	os.WriteFile(path.Join(folder, "1.txt"), []byte("1"), 0777)

	status_output := lib.RunGitCommandWithOutput(folder, "status", "")
	if !strings.Contains(status_output, "modified "+"1.txt") {
		t.Errorf("Status command didn't return the expected output: %s should be modified, got %s", "1.txt", status_output)
	}
}

func TestStatus_DeleteCommitedFile(t *testing.T) {
	folder := lib.GenerateRandomString()
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	one_txt, err := os.Create(path.Join(folder, "1.txt"))
	if err != nil {
		t.Errorf("Can't create the file on path %s", "1.txt")
		return
	}
	one_txt.WriteString("1")

	lib.RunGitCommand(folder, "add", ".")
	lib.RunGitCommand(folder, "commit", "some message")

	one_txt.Close()
	//delete the file so it gets picked up
	err = os.Remove(path.Join(folder, "1.txt"))
	if err != nil {
		t.Errorf("Can't delete the file on path %s", "1.txt")
		return
	}
	status_output := lib.RunGitCommandWithOutput(folder, "status", "")
	if !strings.Contains(status_output, "deleted "+"1.txt") {
		t.Errorf("Status command didn't return the expected output: %s should be deleted, got %s", "1.txt", status_output)
	}
}

func TestStatus_AddNewFile(t *testing.T) {
	folder := lib.GenerateRandomString()
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	one_txt, err := os.Create(path.Join(folder, "1.txt"))
	if err != nil {
		t.Errorf("Can't create the file on path %s", "1.txt")
		return
	}
	one_txt.WriteString("1")

	lib.RunCustomCommand(folder, "mkdir", "a")
	two_txt, err := os.Create(path.Join(folder, "a/2.txt"))
	if err != nil {
		t.Errorf("Can't create the file on path %s", "a/2.txt")
		return
	}
	defer two_txt.Close()
	two_txt.WriteString("2")

	lib.RunGitCommand(folder, "add", ".")

	status_output := lib.RunGitCommandWithOutput(folder, "status", "")
	if !strings.Contains(status_output, "new file "+"1.txt") {
		t.Errorf("Status command didn't return the expected output: %s should be added, got %s", "1.txt", status_output)
	}
	if !strings.Contains(status_output, "new file "+"a/2.txt") {
		t.Errorf("Status command didn't return the expected output: %s should be added, got %s", "2.txt", status_output)
	}
}
