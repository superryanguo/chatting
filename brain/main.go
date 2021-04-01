package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"brain/handler"
	"brain/subscriber"

	brain "brain/proto/brain"
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
