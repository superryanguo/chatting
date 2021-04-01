package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"web-srv/handler"
	"web-srv/subscriber"

	web "web-srv/proto/web"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("micro.chatting.service.web"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	web.RegisterWebHandler(service.Server(), new(handler.Web))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("micro.chatting.service.web", service.Server(), new(subscriber.Web))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
