package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	brain "github.com/superryanguo/chatting/brain/proto/brain"
)

type Brain struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Brain) Call(ctx context.Context, req *brain.Request, rsp *brain.Response) error {
	log.Info("Received Brain.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Brain) Stream(ctx context.Context, req *brain.StreamingRequest, stream brain.Brain_StreamStream) error {
	log.Infof("Received Brain.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&brain.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Brain) PingPong(ctx context.Context, stream brain.Brain_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&brain.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
