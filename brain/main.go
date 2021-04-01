package main

import (
	"github.com/superryanguo/chatting/brain/handler"
	"github.com/superryanguo/chatting/brain/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	brain "github.com/superryanguo/chatting/brain/proto/brain"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("micro.chatting.service.brain"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	brain.RegisterBrainHandler(service.Server(), new(handler.Brain))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("micro.chatting.service.brain", service.Server(), new(subscriber.Brain))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
