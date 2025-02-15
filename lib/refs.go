package lib

import (
	"os"
	"path"
)

type Refs struct {
	Pathname string
}

func NewRefs(pathName string) *Refs {
	return &Refs{Pathname: pathName}
}

func (r *Refs) UpdateHead(oid string) error {
	file, err := os.OpenFile(r.GetHeadPath(), os.O_WRONLY|os.O_CREATE, 0777)
	defer file.Close()
	if err != nil {
		return err
	}
	file.Write([]byte(oid))
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
	return string(b), nil
}
