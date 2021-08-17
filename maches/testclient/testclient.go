package main

import (
	"fmt"
	"time"

	log "github.com/micro/go-micro/v2/logger"

	"github.com/golang/protobuf/proto"
	"github.com/superryanguo/chatting/maches"
	mache "github.com/superryanguo/chatting/websrv/proto/mache"
)

//For test purpose
const (
	Addr = "127.0.0.1:8099"
)

func main() {
	go maches.GetMache().RunClient(Addr)
	var as mache.ChatAnswer
	for i := 0; i < 3; i++ {
		ak := mache.ChatAsk{
			SessionId: "123123",
			Query:     "how are you",
		}
		//maches.GetClientChanRcv() <- ([]byte)("hello")
		data, err := proto.Marshal(&ak)
		if err != nil {
			log.Debug("protoc marshal err=", err.Error())
			return
		}
		log.Debug("sending... the msg:", data)
		maches.GetClientChanRcv() <- data

		select {
		case msg, ok := <-maches.GetClientChanSed():
			if !ok {
				log.Debug("clientChanSed closed")
				return
			}
			log.Debug("RetrieveAnswer receive the msg:", msg)
			err = proto.Unmarshal(msg, &as)
			if err != nil {
				log.Debug("protoc unmarshal err=", err.Error())
				return
			}
			log.Debug("Get the answer=", as.Reply, ", id=", as.SessionId)

		}

		time.Sleep(1 * time.Second)
	}
	fmt.Println("exiting...")
}
