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
消息应答机制确保即使在消费者程序异常未完成消息处理时，消息不会丢失。

在消费者未完成消息处理异常退出时，rabbitmq 会重新将消息收回到队列中。 如果同一个时间存在其他消费者时，rabbitmq 会迅速的将消息发送给其他消费者。

Rabbitmq 在传递(推送)消息过程中会携带 Delivery Tags，用作识别确认消息交付。

Delivery tags 唯一的标识 channel 中的传递，并且作用域为每一个 channel。因此消息确认和接收消息的 channel 要一致，否则会出现协议异常(unknown delivery tag)并关闭 channel

消息应答有两种模式：
- 自动模式——不安全的，不适用与所有的工作负载，并且会导致消费者过载
- 手动模式——手动的有 positive or negative 两种
  - ack positive: 确认消息处理完成，并删除
  - nack + requeue negative: 消息重新排队，再次处理
  - reject negative: 消息路由到死信(dead-letter)队列



