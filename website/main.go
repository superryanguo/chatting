package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"website/handler"
	"website/subscriber"

	website "website/proto/website"
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
