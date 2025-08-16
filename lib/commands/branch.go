package commands

import "geo-git/lib"

func RunBranch(repo *lib.Respository, cmd *Command) error {
	return createBranch(repo, cmd)
}

func createBranch(repo *lib.Respository, cmd *Command) error {
	branchName := cmd.Args[0]
	err := repo.Refs.CreateBranch(branchName)
	return err
}
