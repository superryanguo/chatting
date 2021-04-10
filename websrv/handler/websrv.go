package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	"github.com/superryanguo/chatting/utils"
	websrv "github.com/superryanguo/chatting/websrv/proto/websrv"
)

type Websrv struct{}

func (e *Websrv) Chat(ctx context.Context, req *websrv.ChatRequest, rsp *websrv.ChatResponse) error {
	log.Info("Websrv->Chat func in...")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	rsp.Reply = "hardcode reply..."
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Websrv) Call(ctx context.Context, req *websrv.Request, rsp *websrv.Response) error {
	log.Info("Received Websrv.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Websrv) Stream(ctx context.Context, req *websrv.StreamingRequest, stream websrv.Websrv_StreamStream) error {
	log.Infof("Received Websrv.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&websrv.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Websrv) PingPong(ctx context.Context, stream websrv.Websrv_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&websrv.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
