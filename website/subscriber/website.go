package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	website "github.com/superryanguo/chatting/website/proto/website"
)

type Website struct{}

func (e *Website) Handle(ctx context.Context, msg *website.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *website.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
