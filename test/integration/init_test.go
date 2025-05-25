package integration_test

import (
	lib "geo-git/test/integration/lib"
	"os"
	"path"
	"testing"
)

func TestInit(t *testing.T) {
	folder := lib.GenerateRandomString()
	lib.RunInit(folder)
	defer lib.CleanUpFolder(folder)
	_, err := os.Stat(folder)
	if err != nil {
		t.Errorf("Folder %s does not exist", folder)
	}
	_, err = os.Stat(path.Join(folder, ".git"))
	if err != nil {
		t.Errorf("Folder .git does not exist")
	}
	_, err = os.Stat(path.Join(folder, ".git/objects"))
	if err != nil {
		t.Errorf("Folder objects does not exist")
	}
	_, err = os.Stat(path.Join(folder, ".git/refs"))
	if err != nil {
		t.Errorf("Folder refs does not exist")
	}
}

func TestInit_AlreadyInitExists(t *testing.T) {
	folder := lib.GenerateRandomString()
	lib.RunInit(folder)
	defer func() {
		if r := recover(); r != nil {
			t.Log("Test pased, panic was caught as expected")
		} else {
			t.Errorf("Expected the second call to RunInit() to fail")
		}
		lib.CleanUpFolder(folder)
	}()
	lib.RunInit(folder)
}
