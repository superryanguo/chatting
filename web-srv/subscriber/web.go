package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	web "web-srv/proto/web"
)

type Web struct{}

func (e *Web) Handle(ctx context.Context, msg *web.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *web.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
