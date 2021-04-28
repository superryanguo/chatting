package maches

import (
	"io"
	"net"
	"time"

	log "github.com/micro/go-micro/v2/logger"
)

const (
	MaxConnect = 100
	BufSize    = 1024
)

var (
	mymache *Mache
)

func init() {
	mymache = &Mache{}
	mymache.clientChanSed = make(chan []byte)
	mymache.clientChanRcv = make(chan []byte)
	mymache.serverChanSed = make(chan []byte)
	mymache.serverChanRcv = make(chan []byte)
}

type Mache struct {
	clientChanSed chan []byte
	clientChanRcv chan []byte
	serverChanSed chan []byte
	serverChanRcv chan []byte
}

func GetMache() *Mache {
	return mymache
}

func GetClientChanRcv() chan []byte {
	return mymache.clientChanRcv
}

func GetClientChanSed() chan []byte {
	return mymache.clientChanSed
}

func GetServerChanSed() chan []byte {
	return mymache.serverChanSed
}
func GetServerChanRcv() chan []byte {
	return mymache.serverChanRcv
}

func (m *Mache) RunClient(addr string) {
	var conn net.Conn
	var err error
	for i := 0; i < MaxConnect; i++ {
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			log.Debug("Successfully connect to the addr: ", addr)
			break
		} else {
			log.Info("Fail to connect to the addr: ", addr, "retry...i=", i)
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		log.Fatal("Fail to connect to the MacheServer, err=", err.Error())
	}

	if conn != nil {
		defer conn.Close()
		buf := make([]byte, BufSize)
		for {
			select {
			case message, ok := <-mymache.clientChanRcv:
				if !ok {
					log.Debug("clientChan closed")
					return
				}

				log.Debug("clientChan in message{", message, "}")
				//if len(message) != 0 {
				n, err := conn.Write(message)
				if err != nil {
					log.Info("Write ML endpoint Error:", err.Error())
					break
				}
				log.Debug("maches send out msg to MlEndpoint n=", n)
				n, err = conn.Read(buf)
				if err != nil {
					log.Info("Read ML endpoint Error:", err.Error())
					break
				}
				log.Debug("maches recevie msg from MlEndpointn=", n)
				mymache.clientChanSed <- buf[:n]
				log.Debug("clientChan out message{", buf[:n], "}")
			}
			//}
		}
	}
}

func (m *Mache) RunServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	defer listener.Close()

	log.Debugf("Running@%s......\n", addr)
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
	buf := make([]byte, BufSize)
	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			mymache.serverChanSed <- buf[:n]
			log.Debug("Send to the server chan...", buf[:n])
			msg, ok := <-mymache.serverChanRcv
			if !ok {
				log.Debug("Warning: Channel close")
				return
			}
			log.Debug("Send to the endpoint...", msg)
			_, err = conn.Write(msg)
			if err != nil {
				log.Info("Write ML endpoint Error:", err.Error())
				break
			}
		case io.EOF:
			log.Debugf("Warning: End of data: %s \n", err)
			return
		default:
			log.Debugf("Error: Reading data : %s \n", err)
			return
		}
	}
}
