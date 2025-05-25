package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

func CombineBytes(b1, b2 byte) byte {
	return b1<<4 | b2
}

func PackHexaDecimal(oid string) []byte {
	bytes, _ := hex.DecodeString(oid)
	return bytes
}

func Int32ToBigEndianBytes(num uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Int16ToBigEndianBytes(num uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
