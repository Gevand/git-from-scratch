package database

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/exp/rand"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource((uint64)(time.Now().UnixNano())))

type Database struct {
	Pathname string
}

func NewDatabase(pathname string) *Database {
	return &Database{Pathname: pathname}
}

func (d *Database) StoreBlob(obj *Blob) error {
	content := fmt.Sprintf("%s %v\000%s", obj.Type, len(obj.Data), obj.Data)
	h := sha1.New()
	h.Write([]byte(content))
	obj.Oid = hex.EncodeToString(h.Sum(nil))
	return d.WriteObject(obj.Oid, content)
}

func (d *Database) StoreTree(tree *Tree) error {
	blob := NewBlob(tree.ToString())
	blob.Type = "tree"
	err := d.StoreBlob(blob)
	if err != nil {
		return err
	}
	tree.Oid = blob.Oid
	return nil
}

func (d *Database) StoreCommit(commit *Commit) error {

	blob := NewBlob(commit.ToString())
	blob.Type = "commit"
	err := d.StoreBlob(blob)
	if err != nil {
		return err
	}
	commit.Oid = blob.Oid
	return nil
}

func (d *Database) WriteObject(oid, content string) error {
	object_path := path.Join(d.Pathname, oid[:2], oid[2:])
	dirname := filepath.Dir(object_path)
	temp_path := path.Join(dirname, generateTempName())

	file, err := os.OpenFile(temp_path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0777)
	defer file.Close()
	if err != nil {
		//try to make if path doesn't exist
		if os.IsNotExist(err) {
			err = os.MkdirAll(dirname, 0777)
			if err != nil {
				return err
			}
			file, err = os.OpenFile(temp_path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0777)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	//compress and save
	var buffer bytes.Buffer
	w := zlib.NewWriter(&buffer)
	_, err = w.Write([]byte(content))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, &buffer)
	if err != nil {
		return err
	}

	err = os.Rename(temp_path, object_path)
	return err
}

const charset = "abcdefghijklmnopqrstuvwxyz0987654321"

func generateTempName() string {
	length := 5
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return "temp_obj_" + string(b)
}
