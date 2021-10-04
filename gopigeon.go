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

func HandleConn(c net.Conn) error {
	defer c.Close()
	disconnect := false
	fh, err := mqtt.ReadFixedHeader(c)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	if fh.PktType != mqtt.CONNECT {
		fmt.Println(mqtt.ExpectingConnectPktError)
		return mqtt.ExpectingConnectPktError
	}	

	mqtt.HandleConnect(c, fh)

	for !disconnect {
		fh, err := mqtt.ReadFixedHeader(c)
		// maybe use a switch here or use et method
		if err != nil {
			fmt.Println(err)
			return err
		}	
		switch fh.PktType {
		case mqtt.PUBLISH:
			err = mqtt.HandlePublish(c, fh)
		case mqtt.SUBSCRIBE:
			err = mqtt.HandleSubscribe(c, fh)
		case mqtt.DISCONNECT:
			disconnect = true
		case mqtt.CONNECT:
			// another connect packet was received
			fmt.Println(mqtt.SecondConnectPktError)
			return mqtt.SecondConnectPktError
		default:
			fmt.Println(mqtt.UnknownPktError)
			return mqtt.UnknownPktError
    	}
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}