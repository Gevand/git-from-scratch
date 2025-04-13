package index

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"geo-git/lib/utils"
	"os"
	"strings"
	"syscall"
	"time"
)

const (
	REGULAR_MODE    = "0100644"
	EXECUTABLE_MODE = "0100755"
	MAX_PATH_SIZE   = 0xfff
	ENTRY_BLOCK     = 8
	ENTRY_MIN_SIZE  = 64
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

func ParseEntry(data []byte) (*IndexEntry, error) {
	//N10H40nZ*
	//TODO: Create errors if the bytes passed in can't be parsed
	n := 0
	c_time := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	c_time_nsec := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	m_time := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	m_time_nsec := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	device := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	inode := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	mode := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	uid := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	gid := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	size := binary.BigEndian.Uint32(data[n : n+4])
	n += 4
	oid := hex.EncodeToString(data[n : n+20])
	n += 20
	flags := binary.BigEndian.Uint16(data[n : n+2])
	n += 2
	//need to trim the "padded" 0 bytes at the end
	path := strings.TrimRight(string(data[n:]), "\000")
	return &IndexEntry{
		Oid:        oid,
		Uid:        uid,
		Gid:        gid,
		Size:       int64(size),
		Mode:       mode,
		Inode:      uint64(inode),
		Device:     uint64(device),
		Ctime:      time.Unix(int64(c_time), int64(c_time_nsec)),
		Ctime_Nsec: int64(c_time_nsec),
		Mtime:      time.Unix(int64(m_time), int64(m_time_nsec)),
		Mtime_Nsec: int64(m_time_nsec),
		Flags:      int(flags),
		Path:       path,
	}, nil
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(ts.Sec, ts.Nsec)
}
