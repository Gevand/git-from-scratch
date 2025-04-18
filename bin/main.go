package main

import (
	"bufio"
	"fmt"
	"geo-git/lib/commands"
	"geo-git/lib/database"
	"os"
	"path"
	"time"
)

func main() {
	command := os.Args[1]
	root_path, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "init failed, %v\r\n", err)
		os.Exit(1)
	}
	switch command {
	case "init":
		err = commands.RunInit(root_path, path.Join(root_path, ".git"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "init failed, %v\r\n", err)
			os.Exit(1)
		}
		break
	case "commit":
		name := os.Getenv("GIT_AUTHOR_NAME")
		email := os.Getenv("GIT_AUTHOR_EMAIL")
		if name == "" || email == "" {
			fmt.Fprintf(os.Stderr, "commit failed, %v \r\n", "need an author and email")
			os.Exit(1)
		}
		author := database.NewAuthor(name, email, time.Date(2021, 10, 1, 0, 0, 0, 0, time.Local))
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if message == "" {
			fmt.Fprintf(os.Stderr, "commit failed, %v \r\n", "need a commit message")
			os.Exit(1)
		}
		err = commands.RunCommit(root_path, author, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit failed, %v \r\n", err)
			os.Exit(1)
		}
		break
	case "add":
		err = commands.RunAdd(root_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "add failed, %v \r\n", err)
			os.Exit(1)
		}
		break
	default:
		fmt.Fprintf(os.Stderr, "%v is not a known command\r\n", command)
		os.Exit(1)
	}
}
