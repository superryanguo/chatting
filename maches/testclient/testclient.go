package main

import (
	"github.com/superryanguo/chatting/maches"
)

//For test purpose
const (
	Addr = "127.0.0.1:8099"
)

func main() {
	go maches.GetMache().RunClient(Addr)
	maches.GetClientChanRcv() <- ([]byte)("hello")
}
