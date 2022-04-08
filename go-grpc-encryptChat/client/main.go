package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc/metadata"

	pb "study/go-grpc-encryptCommunications/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world1"
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	name  = flag.String("name", defaultName, "Name to greet")
	mutex sync.Mutex
)

// ConsoleLog 加锁输出，防止乱序或中间插入Print数据
func ConsoleLog(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Printf("\n%s\n", message)
}

// Input 客户端输入
func Input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return ""
		} else {
			panic(err)
		}
	}
	return string(line)
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("name", *name))
	r, err := c.ChatIO(ctx)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 创建一个连接管道
	connected := make(chan bool)
	// wg := sync.WaitGroup{}
	// wg.Add(2)
	// 接收 服务端信息
	go func() {
		for {
			res, err := r.Recv()
			if err != nil {
				fmt.Println(err)
				// wg.Done()
				break
			}
			ConsoleLog(res.Message)
			if res.MessageType == pb.ChatRes_CONNECT_FAILED {
				cancel()
				break
			}
			if res.MessageType == pb.ChatRes_CONNECT_SUCCESS {
				connected <- true
			}
		}
	}()
	// 发送信息给服务端
	go func() {
		<-connected
		var (
			Mes string
			err error
		)
		for {
			Mes = Input(*name + ": ")
			if Mes == "exit" {
				cancel()
				break
			}
			err = r.Send(&pb.ChatReq{
				Message: Mes,
			})
			if err != nil {
				fmt.Println("错误: ", err)
				// wg.Done()
				break
			}
		}
	}()
	<-ctx.Done()
	fmt.Println("Bye")
	// wg.Wait()
}
