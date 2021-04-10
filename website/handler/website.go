package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/superryanguo/chatting/utils"
	website "github.com/superryanguo/chatting/website/proto/website"
)

type Website struct{}

func Init() {
	//userClient = user.NewUserSrvService("micro.super.chatting.service.user_srv", client.DefaultClient)
}

func GetChatMsg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("GetChatMsg-> Chatting message in")
	log.Debugf("httpRequest=%v", r)
	keys, ok := r.URL.Query()["message"]

	if !ok || len(keys[0]) < 1 {
		log.Info("Url Param 'message' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]
	log.Info("GetChatMessage=" + string(key))
}

func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("GetIndex-> html show api/v1.0/chatting/house/index")

	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

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
