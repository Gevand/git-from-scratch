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

func RunGitCommand(folder string, arguments ...string) {
	cmd := exec.Command("geo-git", arguments...)
	cmd.Dir = folder
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func RunGitCommandWithOutput(folder string, arguments ...string) string {
	cmd := exec.Command("geo-git", arguments...)
	cmd.Dir = folder
	out, err := cmd.Output()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

func CleanUpFolder(folder string) {
	err := os.RemoveAll(folder)
	if err != nil {
		panic(err)
	}
}

func RunCustomCommand(folder string, command string, arguments ...string) {
	cmd := exec.Command(command, arguments...)
	cmd.Dir = folder
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
