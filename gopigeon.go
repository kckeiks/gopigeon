package gopigeon

import (
	"fmt"
	"net"
	"os"
	"github.com/kckeiks/gopigeon/mqtt"
	"io"
)

func ListenAndServe() {
	ln, err := net.Listen("tcp", ":8080")
	checkError(err)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// maybe log
			// return CONNACK with error but cant send that if we dont have a conn duh
			continue
		}
		go HandleClient(conn)
	}
}

func HandleClient(c net.Conn) {
	defer c.Close()
	fmt.Printf("net.Conn %d", c)
	fmt.Println(c)
	cs := mqtt.ClientSession{}
	disconnect := false
	for !disconnect {
		fh, err := mqtt.ReadFixedHeader(c)
		// maybe use a switch here or use et method
		if err != nil {
			if err != io.EOF {
				panic(err.Error())
			}
			break
		}
		fmt.Printf("Package fixed-header received: %+v\n", fh)
		if cs.ConnectRcvd && fh.PktType == mqtt.CONNECT {
			break
		}
		switch fh.PktType {
		case mqtt.CONNECT:
			mqtt.HandleConnect(c, fh)
			cs.ConnectRcvd = true
		case mqtt.PUBLISH:
			mqtt.HandlePublish(c, fh)
		case mqtt.SUBSCRIBE:
			mqtt.HandleSubscribe(c, fh)
		case mqtt.DISCONNECT:
			disconnect = true
		default:
			panic("Unknown Control Packet!")
    	}
	}
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}