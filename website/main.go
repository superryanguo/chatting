package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/superryanguo/chatting/website/handler"
	"github.com/superryanguo/chatting/website/subscriber"

	website "github.com/superryanguo/chatting/website/proto/website"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.website"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	website.RegisterWebsiteHandler(service.Server(), new(handler.Website))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.website", service.Server(), new(subscriber.Website))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
