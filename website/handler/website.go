package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/superryanguo/chatting/plugins/psession"
	"github.com/superryanguo/chatting/utils"
	website "github.com/superryanguo/chatting/website/proto/website"
	websrv "github.com/superryanguo/chatting/websrv/proto/websrv"
)

type Website struct{}

var (
	webClient websrv.WebsrvService
)

func Init() {
	webClient = websrv.NewWebsrvService("micro.chatting.service.websrv", client.DefaultClient)
}

func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("GetSession-> /api/v1.0/session in")
	ses, err := psession.GetSession(w, r)
	if err != nil {
		log.Info("Error! no session could be allocated!", err.Error())
		http.Error(w, err.Error(), 502) //TODO: make it like the response in w?
		return
	}
	log.Debug("CtsKey=", ses.Values[psession.CtsKey].(string), "+ses.ID=", ses.ID)
	data := make(map[string]string)
	data["name"] = ses.Values[psession.CtsKey].(string)[:8]
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func GetChatMsg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("GetChatMsg-> Chatting message in")
	//log.Debugf("httpRequest=%v", r)

	//TODO: session
	//userlogin, err := r.Cookie("userlogin")
	//if err != nil {
	//resp := map[string]interface{}{
	//"errno":  utils.RECODE_SESSIONERR,
	//"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
	//}
	//w.Header().Set("Content-Type", "application/json")
	//if err := json.NewEncoder(w).Encode(resp); err != nil {
	//http.Error(w, err.Error(), 503)
	//log.Info(err)
	//return
	//}
	//return
	//}

	keys, ok := r.URL.Query()["cmsg"]

	if !ok || len(keys[0]) < 1 {
		log.Info("Url Param 'cmsg' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]
	log.Debug("cmsg=" + string(key))

	ses, err := psession.GetSession(w, r)
	if err != nil {
		log.Info("Error! no session could be allocated!", err.Error())
		http.Error(w, err.Error(), 502)
		return
	}
	log.Debug("CtsKey=", ses.Values[psession.CtsKey].(string), "+ses.ID=", ses.ID)

	rsp, err := webClient.Chat(context.TODO(), &websrv.ChatRequest{
		SessionId: ses.Values[psession.CtsKey].(string),
		Text:      string(key),
	})

	if err != nil {
		log.Debug("websrv calling error! ", err.Error())
		http.Error(w, err.Error(), 502)
		return
	}
	log.Debugf("rsp=%v", rsp)
	data := make(map[string]interface{})
	data["reply"] = rsp.Reply
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		log.Debug("json calling error! ", err.Error())
		return
	}
	return
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
