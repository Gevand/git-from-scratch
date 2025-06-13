package lib

import (
	"os"
	"path"
	"strings"
)

type Refs struct {
	Pathname string
}

func NewRefs(pathName string) *Refs {
	return &Refs{Pathname: pathName}
}

func (r *Refs) UpdateHead(oid string) error {
	lockfile := NewLockFile(r.GetHeadPath())
	err := lockfile.HoldForUpdate()
	if err != nil {
		return err
	}
	lockfile.Write([]byte(oid))
	lockfile.Write([]byte("\n"))
	err = lockfile.Commit()
	if err != nil {
		return err
	}
	return nil
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
