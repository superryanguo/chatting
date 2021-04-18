package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/cli"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
	"github.com/superryanguo/chatting/basic"
	"github.com/superryanguo/chatting/website/handler"
	//"github.com/superryanguo/chatting/models"
)

const (
	webPort = ":8083"
)

func main() {
	basic.Init()
	//models.Init()

	service := web.NewService(
		web.Name("go.micro.service.website"),
		web.Version("latest"),
		web.Address(webPort),
	)
	// Initialise service
	if err := service.Init(
		web.Action(
			func(c *cli.Context) {
				handler.Init()
			}),
	); err != nil {
		log.Fatal(err)
	}

	rou := httprouter.New()
	rou.NotFound = http.FileServer(http.Dir("html"))
	rou.GET("/api/v1.0/chat", handler.GetChatMsg)
	rou.GET("/api/v1.0/session", handler.GetSession)
	service.Handle("/", rou)
	// Register Handler
	//website.RegisterWebsiteHandler(service.Server(), new(handler.Website))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
