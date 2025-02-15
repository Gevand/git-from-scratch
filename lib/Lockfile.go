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
	Lock     any
}

func NewLockFile(file_path, lock_path string) *LockFile {
	return &LockFile{FilePath: file_path, LockPath: lock_path, Lock: nil}
}

func (l *LockFile) HoldForUpdate() error {
	if l.Lock != nil {
		temp, err := os.OpenFile(l.LockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0777)
		if err != nil {
			return err
		}
		l.Lock = temp

	}
	return nil
}
