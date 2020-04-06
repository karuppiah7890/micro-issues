package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-plugins/registry/etcd"
)

// Example is an example
type Example struct{}

// Handler is a handler
func (e *Example) Handler(ctx context.Context, v interface{}) error {
	fmt.Println("v = ", v)
	return nil
}

func main() {

	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	rabbitmq.DefaultRabbitURL = "amqp://guest:guest@127.0.0.1:5672"
	brkrSub := broker.NewSubscribeOptions(
		broker.Queue("queue.default"),
		broker.DisableAutoAck(),
		rabbitmq.DurableQueue(),
	)

	b := rabbitmq.NewBroker(
		rabbitmq.ExchangeName("ExchangeName"),
	)
	b.Init()
	b.Connect()
	s := server.NewServer(server.Broker(b))

	service := micro.NewService(
		micro.Server(s),
		micro.Broker(b),
		micro.Registry(etcdRegistry),
	)

	//Register a subscriber
	micro.RegisterSubscriber(
		"topic",
		service.Server(),
		new(Example).Handler,
		server.SubscriberContext(brkrSub.Context),
		server.SubscriberQueue("queue.default"),
	)

	service.Init()

	if err := service.Run(); err != nil {
		panic(err)
	}
}
