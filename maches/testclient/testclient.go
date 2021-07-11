package main

import (
	"time"

	"github.com/superryanguo/chatting/maches"
)

//For test purpose
const (
	Addr = "127.0.0.1:8099"
)

func main() {
	go maches.GetMache().RunClient(Addr)
	maches.GetClientChanRcv() <- ([]byte)("hello")
	time.Sleep(5 * time.Second)
}
