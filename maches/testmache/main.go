package main

import (
	"io"
	"net"

	"github.com/golang/protobuf/proto"
	log "github.com/micro/go-micro/v2/logger"
	mache "github.com/superryanguo/chatting/websrv/proto/mache"
)

//For test purpose
const (
	Addr    = "127.0.0.1:8099"
	BufSize = 1024
)

func main() {

	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		panic(err.Error())
	}
	defer listener.Close()

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
	//var buf bytes.Buffer
	buf := make([]byte, BufSize)
	var ak mache.ChatAsk
	defer conn.Close()

	for {
		n, err := conn.Read(buf) //TODO: why it become non-blocking type read? WTF? Because the bytes.Buffer?!
		//n, err := conn.Read(buf.Bytes()) //TODO: why it become non-blocking type read? WTF?
		switch err {
		case nil:
			log.Infof("Testmache Server receive n=%d msg...%v", n, buf[:n])
			err = proto.Unmarshal(buf[:n], &ak)
			if err != nil {
				log.Debug("protoc unmarshal err=", err.Error())
				return
			}
			log.Debugf("Unmarshal data=%v\n", ak)
			as := mache.ChatAnswer{
				SessionId: "1234567890",
				Reply:     "answer from ML",
			}
			data, err := proto.Marshal(&as)
			if err != nil {
				log.Debug("protoc marshal err=", err.Error())
				return
			}
			n, err = conn.Write(data)
			if err != nil {
				log.Info("Write ML endpoint Error:", err.Error())
				return
			}
			log.Infof("Testmache Server Send n=%d msg...%v", n, data)
		case io.EOF:
			log.Debugf("Warning: End of data: %s \n", err)
			return
		default:
			log.Debugf("Error: Reading data : %s \n", err)
			return
		}
	}
}
