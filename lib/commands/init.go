package commands

import (
	"fmt"
	"os"
	"path"
)

var dirs = []string{"objects", "refs"}

func RunInit(root_path, git_path string) {
	err := os.Mkdir(git_path, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatals: %v", err)
		os.Exit(1)
	}
	for _, dir := range dirs {
		err := os.Mkdir(path.Join(git_path, dir), 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatals: %v", err)
			os.Exit(1)
		}
	}
	fmt.Printf("Initialized empty repository in %s\r\n", git_path)
}
