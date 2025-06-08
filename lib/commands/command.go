package commands

import (
	"errors"
	"fmt"
	"geo-git/lib"
	"os"
	"path"
)

type Command struct {
	Args    []string
	EnvVars map[string]string
}

func NewCommand(args []string, env_vars map[string]string) *Command {

	return &Command{Args: args, EnvVars: env_vars}
}

func (c *Command) Execute(name string) error {
	root_path, err := os.Getwd()
	git_path := path.Join(root_path, ".git")
	if err != nil {
		return err
	}
	repo := lib.NewRepository(git_path)
	switch name {
	case "init":
		err = RunInit(repo, c)
	case "commit":
		err = RunCommit(repo, c)
		break
	case "add":
		err = RunAdd(repo, c)
		break
	case "status":
		err = RunStatus(repo, c)
		break
	case "showhead":
		err = RunShowHead(repo)
		break
	default:
		return errors.New(fmt.Sprintf("%s is an unknown command.", name))
	}
	return err
}
