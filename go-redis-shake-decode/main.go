// package main

// // import "C"

// func main() {
// 	var ch chan struct{}
// 	<-ch
// }
package main

import (
	"bufio"
	"fmt"
	utils "go-redis-shake-decode/common"
	"go-redis-shake-decode/libs/atomic2"
)

type CmdDecode struct {
	rbytes atomic2.Int64
}

func main() {
	// 测试 redis-shake dump 模式报错问题
	// 2022/04/23 09:45:49 [PANIC] parse rdb header error
	// [error]: EOF
	// 	5   github.com/alibaba/RedisShake/pkg/rdb/reader.go:102
	// 			github.com/alibaba/RedisShake/pkg/rdb.(*rdbReader).Read
	// 	4   io/io.go:328
	// 			io.ReadAtLeast
	// 	3   io/io.go:347
	// 			io.ReadFull
	// 	2   github.com/alibaba/RedisShake/pkg/rdb/reader.go:445
	// 			github.com/alibaba/RedisShake/pkg/rdb.(*rdbReader).readFull
	// 	1   github.com/alibaba/RedisShake/pkg/rdb/loader.go:34
	// 			github.com/alibaba/RedisShake/pkg/rdb.(*Loader).Header
	// 	0   github.com/alibaba/RedisShake/redis-shake/common/utils.go:953
	// 			github.com/alibaba/RedisShake/redis-shake/common.NewRDBLoader.func1
	readin, _ := utils.OpenReadFile("static/local_dump.0")
	defer readin.Close()

	// saveto := utils.OpenWriteFile(output)
	// defer saveto.Close()

	reader := bufio.NewReaderSize(readin, 8)
	// writer := bufio.NewWriterSize(saveto, utils.WriterBufferSize)
	cmd := new(CmdDecode)
	ipipe := utils.NewRDBLoader(reader, &cmd.rbytes, 1024)
	fmt.Println(ipipe)
	// time.Sleep(10 * time.Second)
	// opipe := make(chan string, cap(ipipe))

}
