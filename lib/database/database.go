package database

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource((uint64)(time.Now().UnixNano())))

type Database struct {
	Pathname string
	Objects  map[string]*Blob
}

func NewDatabase(pathname string) *Database {
	return &Database{Pathname: pathname, Objects: map[string]*Blob{}}
}

func (d *Database) Load(oid string) error {
	blob, err := d.ReadObject(oid)
	d.Objects[oid] = blob
	return err
}

func (d *Database) ReadObject(oid string) (*Blob, error) {
	fmt.Println("Working on ", oid)
	objectPath := d.ObjectPath(oid)
	objectRawText, err := os.ReadFile(objectPath)
	if err != nil {
		fmt.Println("Readfile panic")
		return nil, err
	}

	compressedBuffer := bytes.NewBuffer(objectRawText)
	zlibReader, err := zlib.NewReader(compressedBuffer)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	objectMetaData, err := io.ReadAll(zlibReader)
	if err != nil {
		return nil, err
	}
	objectData := ""
	objectType := ""
	split := strings.Split(string(objectMetaData), "\000")
	objectTypeAndLength := split[0]
	objectType = strings.Split(objectTypeAndLength, " ")[0]
	objectData = strings.Replace(string(objectMetaData), objectType+"\000", "", 1)
	blobToReturn := &Blob{Data: []byte(objectData), Type: objectType}
	blobToReturn.Oid = blobToReturn.HashObject()
	return blobToReturn, nil
}

func (d *Database) StoreBlob(obj *Blob) error {
	content := obj.GetContent()
	obj.Oid = obj.HashObject()
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

func (d *Database) ObjectPath(oid string) string {
	object_path := path.Join(d.Pathname, oid[:2], oid[2:])
	return object_path
}
func (d *Database) WriteObject(oid, content string) error {
	object_path := d.ObjectPath(oid)
	_, err := os.Stat(object_path)
	if !os.IsNotExist(err) {
		return nil
	}
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
	err = w.Close()
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
