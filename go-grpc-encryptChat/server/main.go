package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/peer"

	pb "study/go-grpc-encryptCommunications/proto"

	"google.golang.org/grpc"
)

// server is used to implement  chat.ChatServer.
type server struct {
	pb.UnimplementedChatServer
}

// ConnectPool 定义一个类，继承异步带锁字典，存入grpc stream对象{ name: stream<obj> }
type ConnectPool struct {
	sync.Map
}

var (
	port                      = flag.Int("port", 50051, "The server port")
	connect_pool *ConnectPool = new(ConnectPool)
)

// Get 类ConnectPool方法
func (p *ConnectPool) Get(name string) pb.Chat_ChatIOServer {
	if stream, ok := p.Load(name); ok {
		return stream.(pb.Chat_ChatIOServer)
	} else {
		return nil
	}
}

func (p *ConnectPool) Add(name string, stream pb.Chat_ChatIOServer) {
	p.Store(name, stream)
}

func (p *ConnectPool) Del(name string) {
	p.Delete(name)
}

func (p *ConnectPool) BroadCast(from, message string) {
	log.Printf("BroadCast from: %s, message: %s\n", from, message)
	p.Range(func(username_i, stream_i interface{}) bool {
		username := username_i.(string)
		stream := stream_i.(pb.Chat_ChatIOServer)
		if username == from {
			return true
		} else {
			stream.Send(&pb.ChatRes{
				Message:     message,
				MessageType: pb.ChatRes_NORMAL_MESSAGE,
			})
		}
		return true
	})
}

// ChatIO Streaming incoming outgoing
func (s *server) ChatIO(stream pb.Chat_ChatIOServer) error {
	peer, _ := peer.FromContext(stream.Context())
	log.Printf("Received new connection. %s", peer.Addr.String())
	md, _ := metadata.FromIncomingContext(stream.Context())
	username := md["name"][0] // 从metadata获取用户名信息，可以理解为请求头的数据
	if connect_pool.Get(username) != nil {
		stream.Send(&pb.ChatRes{
			Message:     fmt.Sprintf("username %s already exists!", username),
			MessageType: pb.ChatRes_CONNECT_FAILED, // 重名
		})
		return nil
	} else {
		connect_pool.Add(username, stream)
		stream.Send(&pb.ChatRes{
			Message:     fmt.Sprintf("%s enter room!", username),
			MessageType: pb.ChatRes_CONNECT_SUCCESS,
		})
	}

	go func() {
		// 阻塞，等待断开连接时触发
		<-stream.Context().Done()
		connect_pool.Del(username)
		connect_pool.BroadCast(username, fmt.Sprintf("%s leval room", username))
	}()

	// 广播，用户进入聊天室
	connect_pool.BroadCast(username, fmt.Sprintf("\n---- %s ----\nWelcome %s!", time.Now(), username))
	// 阻塞接收，该用户后续传来的消息
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		connect_pool.BroadCast(username, fmt.Sprintf("%s: %s", username, req.Message))
	}
	// return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
