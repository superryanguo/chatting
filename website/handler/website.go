package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	website "github.com/superryanguo/chatting/website/proto/website"
)

type Website struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Website) Call(ctx context.Context, req *website.Request, rsp *website.Response) error {
	log.Info("Received Website.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Website) Stream(ctx context.Context, req *website.StreamingRequest, stream website.Website_StreamStream) error {
	log.Infof("Received Website.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&website.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Website) PingPong(ctx context.Context, stream website.Website_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&website.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
