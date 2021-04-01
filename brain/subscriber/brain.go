package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	brain "brain/proto/brain"
)

type Brain struct{}

func (e *Brain) Handle(ctx context.Context, msg *brain.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *brain.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
