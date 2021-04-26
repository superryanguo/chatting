package main

import (
	"bytes"
	"io"
	"net"

	"github.com/golang/protobuf/proto"
	log "github.com/micro/go-micro/v2/logger"
	mache "github.com/superryanguo/chatting/websrv/proto/mache"
)

//For test purpose
const (
	Addr = "127.0.0.1:9990"
)

func main() {

	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		panic(err.Error())
	}

	log.Debugf("Running@%s......\n", Addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Error accept:" + err.Error())
		}
		log.Debug("Accepted the Connection :", conn.RemoteAddr())
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	var buf bytes.Buffer
	var ak mache.ChatAsk
	defer conn.Close()

	for {
		_, err := conn.Read(buf.Bytes())
		switch err {
		case nil:
			log.Debug("Receive msg...", buf.Bytes())
			err = proto.Unmarshal(buf.Bytes(), &ak)
			if err != nil {
				log.Debug("protoc unmarshal err=", err.Error())
				return
			}
			log.Debugf("Unmarshal data=%v\n", ak)
			as := mache.ChatAnswer{
				SessionId: "1234567890",
				Reply:     "answer",
			}
			data, err := proto.Marshal(&as)
			if err != nil {
				log.Debug("protoc marshal err=", err.Error())
				return
			}
			_, err = conn.Write(data)
			if err != nil {
				log.Info("Write ML endpoint Error:", err.Error())
				return
			}
		case io.EOF:
			log.Debug("Warning: End of data: %s \n", err)
			return
		default:
			log.Debug("Error: Reading data : %s \n", err)
			return
		}
	}
}
