package index

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"hash"
	"os"
)

const CHECKSUM_SIZE = 20

type Checksum struct {
	File   *os.File
	Digest hash.Hash
}

func NewChecksum(file *os.File) *Checksum {
	return &Checksum{
		File: file, Digest: sha1.New(),
	}
}

func (c *Checksum) Read(size int) ([]byte, error) {
	data := make([]byte, size)
	n, err := c.File.Read(data)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, errors.New("Unexpected end-of-file while reading index")
	}

	c.Digest.Write(data)

	return data, nil
}

func (c *Checksum) Verify() error {
	sum := make([]byte, CHECKSUM_SIZE)
	_, err := c.File.Read(sum)
	if err != nil {
		return err
	}
	digest := c.Digest.Sum(nil)
	if !bytes.Equal(sum, digest) {
		return errors.New("Checksum does not match value stored on disk")
	}
	return nil

}
