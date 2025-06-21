package database

import (
	"errors"
	"fmt"
	"strings"
)

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

func ParseCommitFromBlob(blob *Blob) (*Commit, error) {
	if blob.Type != "commit" {
		return nil, errors.New("type must be commit, got " + blob.Type + " instead")
	}

	commitToReturn := &Commit{Oid: blob.Oid}
	lines := strings.Split(string(blob.Data), "\n\n")
	if len(lines) < 2 {
		return nil, errors.New("commit is not in the right format")
	}
	//line \n\n used to separate the message from the rest
	commitToReturn.Message = strings.Join(lines[1:], "\n\n")
	//split the rest of them by line
	lines = strings.Split(lines[0], "\n")
	author := Author{}
	for _, line := range lines {
		if strings.HasPrefix(line, "author") {
			err := author.Parse(strings.Replace(line, "author ", "", 1))
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "tree") {
			commitToReturn.Tree_Oid = strings.Replace(line, "tree ", "", 1)
		} else if strings.HasPrefix(line, "commiter") {
			author.Parse(strings.Replace(line, "author ", "", 1))
		} else if strings.HasPrefix(line, "parent") {
			commitToReturn.Parent = strings.Replace(line, "parent ", "", 1)
		}
	}
	commitToReturn.Author = author
	return commitToReturn, nil
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
		"tree %s\nparent %s\nauthor %s\ncommitter %s\n\n%s",
		c.Tree_Oid,
		c.Parent,
		c.Author.ToString(),
		c.Author.ToString(),
		c.Message)
}
