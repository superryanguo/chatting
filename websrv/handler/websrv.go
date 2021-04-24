package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/micro/go-micro/v2/logger"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/superryanguo/chatting/basic/cache"
	"github.com/superryanguo/chatting/basic/config"
	"github.com/superryanguo/chatting/models"
	"github.com/superryanguo/chatting/utils"
	websrv "github.com/superryanguo/chatting/websrv/proto/websrv"
)

type Websrv struct{}

var (
	dialogPath   string
	dialogPrefix string
)

func Init() {
	dialogPath = config.GetMisconfig().GetDialogPath()
	dialogPrefix = config.GetMisconfig().GetDialogPrefix()
	_ = CreateOneUser()
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
			d.Path = dialogPath + dialogPrefix + uuid.New().String()
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
func (e *Websrv) Chat(ctx context.Context, req *websrv.ChatRequest, rsp *websrv.ChatResponse) error {
	log.Info("Websrv->Chat func in...")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	rsp.Reply = "hardcode reply for ses" + req.SessionId

	path, err := GetSessionProfile(req.SessionId)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	var output []string
	output = append(output, fmt.Sprintf("User: %s\n", req.Text))
	output = append(output, fmt.Sprintf("Robot: %s\n", rsp.Reply))
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	defer f.Close()

	if _, err := f.WriteString(strings.Join(output, "")); err != nil {
		rsp.Errno = utils.RECODE_DATAERR
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
