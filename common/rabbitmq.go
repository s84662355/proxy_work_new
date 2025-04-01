package common

import (
	"sync"

	"mproxy/config"
	"mproxy/utils/rabbitMQ"
)

// 获取rabbitmq生产者
var GetRabbitMqProductPool = sync.OnceValue[*rabbitMQ.RabbitMQPool](func() *rabbitMQ.RabbitMQPool {
	pool := rabbitMQ.NewProductPool()

	err := pool.ConnectVirtualHost(
		config.GetConf().Rabbitmq.Host,
		config.GetConf().Rabbitmq.Port,
		config.GetConf().Rabbitmq.User,
		config.GetConf().Rabbitmq.Password,
		config.GetConf().Rabbitmq.VirtualHost,
	)
	if err != nil {
		panic(err)
	}
	return pool
})

// 获取rabbitmq的消费者
var GetRabbitConsumerPool = sync.OnceValue[*rabbitMQ.RabbitMQPool](func() *rabbitMQ.RabbitMQPool {
	pool := rabbitMQ.NewConsumePool()

	err := pool.ConnectVirtualHost(
		config.GetConf().Rabbitmq.Host,
		config.GetConf().Rabbitmq.Port,
		config.GetConf().Rabbitmq.User,
		config.GetConf().Rabbitmq.Password,
		config.GetConf().Rabbitmq.VirtualHost,
	)
	if err != nil {
		panic(err)
	}
	return pool
})
