package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/micro/go-micro/v2/logger"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/superryanguo/chatting/basic/cache"
	"github.com/superryanguo/chatting/basic/config"
	"github.com/superryanguo/chatting/maches"
	"github.com/superryanguo/chatting/models"
	"github.com/superryanguo/chatting/utils"
	mache "github.com/superryanguo/chatting/websrv/proto/mache"
	websrv "github.com/superryanguo/chatting/websrv/proto/websrv"
)

type Websrv struct{}

var (
	dialogPath   string
	dialogPrefix string
	macheEndAddr string
)

func Init() {
	dialogPath = config.GetMisconfig().GetDialogPath()
	dialogPrefix = config.GetMisconfig().GetDialogPrefix()
	macheEndAddr = config.GetMisconfig().GetMLAddr() + ":"
	macheEndAddr += strconv.Itoa(config.GetMisconfig().GetMLPort())
	_ = CreateOneUser()
	if err := utils.CreateDir(dialogPath); err != nil {
		log.Info("Dir not exist or fail to create", err.Error())
	}
	go maches.GetMache().RunClient(macheEndAddr)
}

func CreateOneUser() error {
	var u models.User
	db := models.GetGorm()
	err := db.Debug().First(&u, models.Userid).Error
	if err != nil {
		log.Debug("CreateOneUser can't find user in DB, Err:", err)
		u.Name = "helloworld"
		u.Email = "hellowrold@chat.com"
		u.RealName = "kick"
		u.AvatarUrl = "127.0.0.1:3300/img111"
		u.ID = models.Userid

		err = db.Debug().Create(&u).Error
		if err != nil {
			log.Debug("CreateOneUser->fail to insert a user to db", u)
			return err
		}
	}
	log.Debug("User in the db done, u:", u)
	return nil
}

//GetSessionProfile get the session's dialog file or create a new one
//if the session is new
func GetSessionProfile(sid string) (path string, err error) {
	var d models.Dialog

	dialog, err := cache.GetFromCache(sid)
	if err != nil {
		if err == redis.Nil {
			log.Debug("GetSessionProfile->no this sid in cache redis.Nil")
			d.SessionId = sid
			d.UserID = models.Userid //TODO: make this one select from the db
			d.Path = dialogPath + dialogPrefix + sid
			log.Debug("GetSessionProfile->create new dialog:", d)
			dj, err := json.Marshal(d)
			if err != nil {
				log.Debug("GetSessionProfile->json problem:", err.Error())
				return "", err
			}
			err = cache.SaveToCache(sid, dj)
			if err != nil {
				log.Debug("GetSessionProfile->redis save failure:", err.Error())
				return "", err
			}
			//save to database
			//TODO: check the dialog exist first is better? then update
			err = models.GetGorm().Debug().Create(&d).Error
			if err != nil {
				log.Debug("GetSessionProfile->fail to insert a dialog to db", err.Error())
				return "", err
			}
			return d.Path, nil
		} else {
			log.Debug("GetSessionProfile->cache problem:", err.Error())
			return "", err
		}
	}

	err = json.Unmarshal(dialog, &d)
	if err != nil {
		log.Debug("GetSessionProfile->json problem:", err.Error())
		return "", err
	}
	return d.Path, nil

}

func RetrieveAnswer(id, ask string) (string, error) {
	var as mache.ChatAnswer
	ak := mache.ChatAsk{
		SessionId: id,
		Query:     ask,
	}
	data, err := proto.Marshal(&ak)
	if err != nil {
		log.Debug("protoc marshal err=", err.Error())
		return "", err
	}
	log.Debug("RetrieveAnswer sending the msg:", data)
	maches.GetClientChanRcv() <- data

	select {
	case msg, ok := <-maches.GetClientChanSed():
		if !ok {
			log.Debug("clientChanSed closed")
			return "", errors.New("clientChanSed closed")
		}
		log.Debug("RetrieveAnswer receive the msg:", msg)
		err = proto.Unmarshal(msg, &as)
		if err != nil {
			log.Debug("protoc unmarshal err=", err.Error())
			return "", err
		}
	}
	return as.Reply, nil
}

func (e *Websrv) Chat(ctx context.Context, req *websrv.ChatRequest, rsp *websrv.ChatResponse) error {
	log.Info("Websrv->Chat func in...")
	var err error
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//TODO: get a reply from ML, hardcode one now
	//rsp.Reply = "hardcode reply for ses" + req.SessionId
	rsp.Reply, err = RetrieveAnswer(req.SessionId, req.Text)
	if err != nil {
		log.Debug("Chat->RetrieveAnswer error:", err.Error())
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	path, err := GetSessionProfile(req.SessionId)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	log.Debug("Get the path=", path)

	var output []string
	output = append(output, fmt.Sprintf("User: %s\n", req.Text))
	output = append(output, fmt.Sprintf("Robot: %s\n", rsp.Reply))
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Debug("Chat->File open error:", err.Error())
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	defer f.Close()

	if _, err := f.WriteString(strings.Join(output, "")); err != nil {
		log.Debug("Chat->File write error:", err.Error())
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

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
