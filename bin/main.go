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
		}
		commands.RunInit(root_path, path.Join(root_path, ".git"))
		break
	default:
		fmt.Fprintf(os.Stderr, "%v is not a known command\r\n", command)
	}
}
