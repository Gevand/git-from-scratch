package lib

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const INVALID_NAME = `/
^\.
| \/\.
| \.\.
| ^\/
| \/$
| \.lock$
| @\{
| [\x00-\x20*:?\[\\^~\x7f]
/x`

type Refs struct {
	Pathname  string
	HeadsPath string
	RefsPath  string
}

func NewRefs(pathName string) *Refs {
	rPath := path.Join(pathName, "refs")
	hPath := path.Join(rPath, "heads")
	return &Refs{Pathname: pathName, HeadsPath: hPath, RefsPath: rPath}
}

func (r *Refs) UpdateHead(oid string) error {
	return r.UpdateRefFile(r.GetHeadPath(), oid, true)
}

func (r *Refs) GetHeadPath() string {
	return path.Join(r.Pathname, "HEAD")
}

func (r *Refs) ReadHead() (string, error) {
	b, err := os.ReadFile(r.GetHeadPath())
	if os.IsNotExist(err) {
		return "", nil
	}

	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

type BranchError struct {
	Message string
}

func (e *BranchError) Error() string {
	return e.Message
}

func (r *Refs) CreateBranch(branchName string) error {
	branchPath := path.Join(r.HeadsPath, branchName)
	matched, err := regexp.MatchString(INVALID_NAME, branchName)
	if matched {
		return &BranchError{Message: "A branch named " + branchName + " has an invalid name"}
	}
	if err != nil {
		return err
	}
	_, err = os.Stat(branchPath)
	if err == nil {
		return &BranchError{Message: "A branch named " + branchName + " alread exists"}
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	head_oid, err := r.ReadHead()
	if err != nil {
		return err
	}
	err = r.UpdateRefFile(branchPath, head_oid, true)
	if err != nil {
		return err
	}
	return nil
}

func (r *Refs) UpdateRefFile(path, oid string, retry bool) error {
	lockfile := NewLockFile(path)
	err := lockfile.HoldForUpdate()
	if retry && errors.Is(err, os.ErrNotExist) {
		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, 0777)
		r.UpdateRefFile(path, oid, false)

	}
	if err != nil {
		return err
	}
	lockfile.Write([]byte(oid))
	lockfile.Write([]byte("\n"))
	lockfile.Commit()

	return nil
}
