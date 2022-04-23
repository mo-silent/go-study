package utils

import (
	"bufio"
	"fmt"
	"go-redis-shake-decode/libs/atomic2"
	"go-redis-shake-decode/libs/stats"
	"go-redis-shake-decode/rdb"
	"os"
)

func OpenReadFile(name string) (*os.File, int64) {
	f, err := os.Open(name)
	if err != nil {
		fmt.Printf("%s cannot open file-reader '%s'", err, name)
	}
	s, err := f.Stat()
	if err != nil {
		fmt.Printf("%s cannot stat file-reader '%s'", err, name)
	}
	fmt.Printf("%v\n%v", f, s.Size())
	return f, s.Size()
}

func NewRDBLoader(reader *bufio.Reader, rbytes *atomic2.Int64, size int) chan *rdb.BinEntry {
	pipe := make(chan *rdb.BinEntry, size)
	go func() {
		defer close(pipe)
		l := rdb.NewLoader(stats.NewCountReader(reader, rbytes))
		if err := l.Header(); err != nil {
			fmt.Printf("%s parse rdb header error", err)
		}
	}()
	return pipe
}
