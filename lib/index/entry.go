package index

import (
	"errors"
	"os"
	"syscall"
	"time"
)

const (
	REGULAR_MODE    = "0100644"
	EXECUTABLE_MODE = "0100755"
	MAX_PATH_SIZE   = 0xfff
)

type IndexEntry struct {
	Ctime, Mtime                 time.Time
	Ctime_Nsec, Mtime_Nsec, Size int64
	Device, Inode                uint64
	Uid, Gid, Mode               uint32
	Oid, Path                    string
	Flags                        int
}

func NewEntry(stat os.FileInfo, path, oid string) (*IndexEntry, error) {
	if stats, ok := stat.Sys().(*syscall.Stat_t); ok {
		return &IndexEntry{
			Uid:        stats.Uid,
			Gid:        stats.Gid,
			Size:       stats.Size,
			Ctime:      timespecToTime(stats.Ctim),
			Mtime:      timespecToTime(stats.Mtim),
			Ctime_Nsec: stats.Ctim.Nsec,
			Mtime_Nsec: stats.Mtim.Nsec,
			Device:     stats.Dev, Inode: stats.Ino,
			Mode: stats.Mode, Path: path, Oid: oid,
			Flags: min(len([]byte(path)), MAX_PATH_SIZE)}, nil
	} else {
		return nil, errors.New("Unable to parse fstat")
	}
}

func (ie *IndexEntry) ToString() string {
	return "TODO"
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(ts.Sec, ts.Nsec)
}
