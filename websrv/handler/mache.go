package handler

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

func Dial() {
	conn, err := net.Dial("tcp", "127.0.0.1:6666")
	defer conn.Close()
	//maintain a connection, setup need 10s*500times to setup the connection
	msg := fmt.Sprintf("Hello World, %03d", i)
	n, err := conn.Write([]byte(msg))
	n, err = conn.Read(buf)
	if err != nil {
		println("Write Buffer Error:", err.Error())
		break
	}
}
func read2(conn *net.Conn) error {
	defer conn.Close()

	var buf bytes.Buffer

	_, err := io.Copy(&buf, conn)
	if err != nil {
		// Error handler
		return err
	}

	return nil
}
func check() {
	listener, err := net.Listen("tcp", "0.0.0.0:6666")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Error accept:" + err.Error())
		}
		fmt.Println("Accepted the Connection :", conn.RemoteAddr())
		go EchoServer(conn)
	}
	for {
		timeout_cnt := 0
		select {
		case msg1 := <-c1:
			fmt.Println("msg1 received", msg1)
		case msg2 := <-c2:
			fmt.Println("msg2 received", msg2)
		case <-time.After(time.Second * 30):
			fmt.Println("Time Out")
			timout_cnt++
		}
		if time_cnt > 3 {
			break
		}
	}
}
func EchoServer(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		switch err {
		case nil:
			conn.Write(buf[0:n])
		case io.EOF:
			fmt.Printf("Warning: End of data: %s \n", err)
			return
		default:
			fmt.Printf("Error: Reading data : %s \n", err)
			return
		}
	}
}
