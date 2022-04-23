package rdb

import (
	"bytes"
	"errors"
	"fmt"
	"go-redis-shake-decode/digest"
	"hash"
	"io"
	"strconv"
)

type Loader struct {
	*rdbReader
	crc       hash.Hash64
	db        uint32
	lastEntry *BinEntry
}

type BinEntry struct {
	DB              uint32
	Key             []byte
	Type            byte
	Value           []byte
	ExpireAt        uint64
	RealMemberCount uint32
	NeedReadLen     byte
	IdleTime        uint32
	Freq            uint8
}

type rdbReader struct {
	raw            io.Reader
	buf            [8]byte
	nread          int64
	remainMember   uint32
	lastReadCount  uint32
	totMemberCount uint32
}

func (l *Loader) Header() error {
	header := make([]byte, 9)
	if err := l.readFull(header); err != nil {
		return err
	}
	if !bytes.Equal(header[:5], []byte("REDIS")) {
		return errors.New("verify magic string, invalid file format")
	}
	if version, err := strconv.ParseInt(string(header[5:]), 10, 64); err != nil {
		return err
	} else if version <= 0 {
		return fmt.Errorf("verify version, invalid RDB version number: %v", fmt.Sprint(version))
	}
	return nil
}

func (r *rdbReader) readFull(p []byte) error {
	_, err := io.ReadFull(r, p)
	return err
}

func (r *rdbReader) Read(p []byte) (int, error) {
	n, err := r.raw.Read(p)
	r.nread += int64(n)
	return n, err
}

func NewLoader(r io.Reader) *Loader {
	l := &Loader{}
	l.crc = digest.New()
	l.rdbReader = NewRdbReader(io.TeeReader(r, l.crc))
	return l
}

func NewRdbReader(r io.Reader) *rdbReader {
	return &rdbReader{raw: r, remainMember: 0, lastReadCount: 0}
}
