package database

import "fmt"

type Commit struct {
	Tree_Oid string
	Author   Author
	Message  string
	Oid      string
}

func NewCommit(tree_oid string, author Author, message string) *Commit {
	return &Commit{Tree_Oid: tree_oid, Author: author, Message: message}
}

func (c *Commit) ToString() string {
	return_value := fmt.Sprintf(
		"tree %s\nauthor %s\ncommitter %s\n\n%s",
		c.Tree_Oid,
		c.Author.ToString(),
		c.Author.ToString(),
		c.Message)
	return return_value
}
