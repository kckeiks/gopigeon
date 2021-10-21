package gopigeon

import (
	"fmt"
	"time"
	"net"
	"testing"
)

func disconnect(c net.Conn, done chan bool){
	c.Close()
	done <- true
}

func server(done chan bool) {
	ln, err := net.Listen("tcp", "localhost:8033")
	if err != nil {
		panic(fmt.Sprintf("Fatal error: %v", err))
	}
	conn, err := ln.Accept()
	if err != nil {
		panic(fmt.Sprintf("Fatal error: %v", err))
	}
	defer disconnect(conn, done)
	b := make([]byte, 5)
	_, err = conn.Read(b)
	if err != nil {
		panic(fmt.Sprintf("Fatal error: %v", err))
	}
	fmt.Printf("Server received: %s\n", b)
}

func client(done chan bool) {
	conn, err := net.Dial("tcp", "localhost:8033")
	if err != nil {
		panic(fmt.Sprintf("Failed to dial: %v\n", err))
	}
	defer disconnect(conn, done)
	_, err = conn.Write([]byte("Hello"))
	if err != nil {
		panic(fmt.Sprintf("dial: error: %v\n", err))
	}
}

func TestPingReq(t *testing.T) {
	done := make(chan bool, 2)
	go server(done)
	time.Sleep(time.Second)
	go client(done)
	<-done
	<-done
}