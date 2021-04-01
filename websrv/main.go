package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/superryanguo/chatting/websrv/handler"
	"github.com/superryanguo/chatting/websrv/subscriber"

	websrv "github.com/superryanguo/chatting/websrv/proto/websrv"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("micro.chatting.service.websrv"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	websrv.RegisterWebsrvHandler(service.Server(), new(handler.Websrv))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("micro.chatting.service.websrv", service.Server(), new(subscriber.Websrv))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
