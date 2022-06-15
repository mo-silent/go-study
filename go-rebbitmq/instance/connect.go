package instance

import "github.com/streadway/amqp"

const URI = "amqp://39.101.244.245:9672/"

type RebbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
}
