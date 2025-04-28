package main

import (
	"fmt"
	"geo-git/lib/commands"
	"os"
)

var env_variables = [...]string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL"}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Program requires arguments to be passed into it \r\n")
		os.Exit(-1)
	}

	command := os.Args[1]
	temp := make(map[string]string)
	for _, env_var := range env_variables {
		temp[env_var] = os.Getenv(env_var)
	}
	cmd := commands.NewCommand(os.Args[2:], temp)
	err := cmd.Execute(command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v \r\n", err)
		os.Exit(-1)
	}
}
