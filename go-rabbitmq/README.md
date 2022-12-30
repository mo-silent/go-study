# RabbitMQ 学习笔记
Git 仓库地址: https://gitee.com/MoGD/go-study/tree/main/go-rabbitmq
## 安装
```shell
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management
```
## 官方七大示例
1. Hello World
2. Work queues
3. Publish/Subscribe
4. Routing
5. Topics
6. RPC
7. Publisher Confirms
### 简单模式-Hello Word
简单模式是一个生产者和一个消费者对应一个队列，一般在生产上不会使用该模式。
```shell
# run producer
go run ./hello-world/send.go
# run consumer
go run ./hello-world/receive.go
```

#### Message Acknowledgment


