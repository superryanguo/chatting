package maches

import (
	"bytes"
	"io"
	"net"
	"time"

	log "github.com/micro/go-micro/v2/logger"
)

const (
	MaxConnect = 100
)

var (
	mymache *Mache
)

func init() {
	mymache.clientChan = make(chan []byte)
	mymache.serverChan = make(chan []byte)
}

type Mache struct {
	clientChan chan []byte
	serverChan chan []byte
}

func GetMache() *Mache {
	return mymache
}

func GetClientChan() chan []byte {
	return mymache.clientChan
}

func GetServerChan() chan []byte {
	return mymache.serverChan
}

func (m *Mache) RunClient(addr string) {
	var conn net.Conn
	var err error
	log.Debug("Addr is ", addr)
	for i := 0; i < MaxConnect; i++ {
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			log.Debug("Successfully connect to the addr: ", addr)
			break
		} else {
			log.Info("Fial to connect to the addr: ", addr, "retry...i=", i)
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		log.Fatal("Fail to connect to the MacheServer, err=", err.Error())
	}

	if conn != nil {
		defer conn.Close()
		var buf bytes.Buffer
		for {
			select {
			case message := <-mymache.clientChan:
				log.Debug("message ", message, " in...")
				_, err := conn.Write(message)
				if err != nil {
					log.Info("Write ML endpoint Error:", err.Error())
					break
				}
				_, err = conn.Read(buf.Bytes())
				if err != nil {
					log.Info("Read ML endpoint Error:", err.Error())
					break
				}
				mymache.clientChan <- buf.Bytes()
				log.Debug("reply msg ", buf.Bytes(), " out...")
			}
		}
	}
}

func (m *Mache) RunServer(addr string) {
	log.Debug("Addr is ", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

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
	defer conn.Close()

	for {
		_, err := conn.Read(buf.Bytes())
		switch err {
		case nil:
			mymache.serverChan <- buf.Bytes()
			log.Debug("Send to the server chan...", buf.Bytes())
			msg := <-mymache.serverChan
			log.Debug("Send to the endpoint...", msg)
			_, err = conn.Write(msg)
			if err != nil {
				log.Info("Write ML endpoint Error:", err.Error())
				break
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
