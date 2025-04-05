package integration_test_lib

import (
	"os"
	"os/exec"
)

func RunInit(folder string) {
	err := os.Mkdir(folder, 777)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("geo-git", "init")
	cmd.Dir = folder
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func CleanUpFolder(folder string) {
	err := os.RemoveAll(folder)
	if err != nil {
		panic(err)
	}
}
