package gopigeon

import (
	"fmt"
	"net"
	"os"
	"github.com/kckeiks/gopigeon/mqtt"
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
		go HandleConn(conn)
	}
}

func HandleConn(c net.Conn) {
	defer c.Close()
	disconnect := false
	fh, err := mqtt.ReadFixedHeader(c)
	if err != nil {
		fmt.Println("error: reading fixed header")
		return
	}
	
	if fh.PktType != mqtt.CONNECT {
		// first pkt was not connect
		fmt.Println("error: first pkt was not CONNECT")
		return 
	}	

	mqtt.HandleConnect(c, fh)

	for !disconnect {
		fh, err := mqtt.ReadFixedHeader(c)
		fmt.Printf("FixedHeader: %+v\n", fh)
		// maybe use a switch here or use et method
		if err != nil {
			fmt.Println("error: reading fixed header")
			return
		}	
		switch fh.PktType {
		case mqtt.PUBLISH:
			mqtt.HandlePublish(c, fh)
		case mqtt.SUBSCRIBE:
			mqtt.HandleSubscribe(c, fh)
		case mqtt.DISCONNECT:
			disconnect = true
		case mqtt.CONNECT:
			// another connect packet was received
			fmt.Println("error: a second CONNECT was received")
			break
		default:
			fmt.Println("error: unknown packet.")
			break
    	}
	}
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}