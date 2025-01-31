package main

import (
	"fmt"
	"geo-git/lib/commands"
	"os"
	"path"
)

func main() {
	command := os.Args[1]
	switch command {
	case "init":
		root_path, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "init failed, %v\r\n", err)
			os.Exit(1)
		}
		commands.RunInit(root_path, path.Join(root_path, ".git"))
		break
	case "commit":
		root_path, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit failed, %v\r\n", err)
			os.Exit(1)
		}
		commands.RunCommit(root_path)

	default:
		fmt.Fprintf(os.Stderr, "%v is not a known command\r\n", command)
		os.Exit(1)
	}
}
