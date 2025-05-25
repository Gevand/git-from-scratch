package lib

import (
	"errors"
	"os"
)

var MissingParent = errors.New("missing parent")
var NoPermission = errors.New("no permissions")
var StaleLock = errors.New("stale lock")

type LockFile struct {
	FilePath string
	LockPath string
	Lock     *os.File
}

func NewLockFile(lock_path string) *LockFile {
	return &LockFile{FilePath: lock_path, LockPath: lock_path + ".lock", Lock: nil}
}

func (l *LockFile) HoldForUpdate() error {
	if l.Lock == nil {
		temp, err := os.OpenFile(l.LockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0777)
		if err != nil {
			return err
		}
		l.Lock = temp
	}
	return nil
}

func (l *LockFile) Write(data []byte) error {
	if l.Lock != nil {
		_, err := l.Lock.Write(data)
		return err
	}
	return errors.New("File needs to be locked before writing to it")
}

func (l *LockFile) Commit() error {
	if l.Lock != nil {
		err := l.Lock.Close()
		if err != nil {
			return err
		}
		err = os.Rename(l.LockPath, l.FilePath)
		l.Lock = nil
		return err
	}
	return errors.New("File needs to be locked before commiting it to disc")

}
