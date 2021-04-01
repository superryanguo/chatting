package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	websrv "github.com/superryanguo/chatting/websrv/proto/websrv"
)

type Websrv struct{}

func (e *Websrv) Handle(ctx context.Context, msg *websrv.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *websrv.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
