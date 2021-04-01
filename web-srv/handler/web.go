package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	web "web-srv/proto/web"
)

type Web struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Web) Call(ctx context.Context, req *web.Request, rsp *web.Response) error {
	log.Info("Received Web.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Web) Stream(ctx context.Context, req *web.StreamingRequest, stream web.Web_StreamStream) error {
	log.Infof("Received Web.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&web.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Web) PingPong(ctx context.Context, stream web.Web_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&web.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
