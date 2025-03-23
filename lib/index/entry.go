package index

import (
	"errors"
	"fmt"
	"geo-git/lib/utils"
	"os"
	"syscall"
	"time"
)

const (
	REGULAR_MODE    = "0100644"
	EXECUTABLE_MODE = "0100755"
	MAX_PATH_SIZE   = 0xfff
	ENTRY_BLOCK     = 8
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

func (ie *IndexEntry) ToString() (string, error) {
	//N10H40nZ*

	//N10 -> 10 32 bit ints int his order -> https://git-scm.com/docs/index-format
	result, err := utils.Int32ToBigEndianBytes(uint32(ie.Ctime.Unix()))
	if err != nil {
		return "", err
	}

	temp, err := utils.Int32ToBigEndianBytes(uint32(ie.Ctime_Nsec))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Mtime.Unix()))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Mtime_Nsec))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Device))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Inode))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Mode))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Uid))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Gid))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	temp, err = utils.Int32ToBigEndianBytes(uint32(ie.Size))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	//H20 -> OID packed from 40 bytes to 20
	temp = utils.PackHexaDecimal(ie.Oid)
	fmt.Println("OID", ie.Oid, len(ie.Oid))
	result = append(result, temp...)

	//n -> 16 bit int
	temp, err = utils.Int16ToBigEndianBytes(uint16(ie.Flags))
	if err != nil {
		return "", err
	}
	result = append(result, temp...)

	//null terimanted result with 0s till its divisible by ENTRY_BLOCK size
	result = append(result, []byte(ie.Path+"\000")...)
	for {
		if len(result)%ENTRY_BLOCK == 0 {
			break
		}
		result = append(result, byte(0))
	}

	return string(result), nil
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(ts.Sec, ts.Nsec)
}
