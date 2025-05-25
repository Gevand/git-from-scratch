package database

import "fmt"

type Commit struct {
	Tree_Oid string
	Author   Author
	Message  string
	Oid      string
	Parent   string
}

func NewCommit(parent string, tree_oid string, author Author, message string) *Commit {
	return &Commit{Tree_Oid: tree_oid, Author: author, Message: message, Parent: parent}
}

func (c *Commit) ToString() string {
	if c.Parent == "" {
		return fmt.Sprintf(
			"tree %s\nauthor %s\ncommitter %s\n\n%s",
			c.Tree_Oid,
			c.Author.ToString(),
			c.Author.ToString(),
			c.Message)
	}
	return fmt.Sprintf(
		"tree %s\nprent %s\nauthor %s\ncommitter %s\n\n%s",
		c.Tree_Oid,
		c.Parent,
		c.Author.ToString(),
		c.Author.ToString(),
		c.Message)
}
