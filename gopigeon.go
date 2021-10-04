package gopigeon

import (
	"fmt"
	"net"
	"os"
)

func ListenAndServe() {
	ln, err := net.Listen("tcp", ":8080")
	checkError(err)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go HandleConn(conn)
	}
}

func HandleConn(c net.Conn) error {
	defer c.Close()
	disconnect := false
	fh, err := ReadFixedHeader(c)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if fh.PktType != CONNECT {
		fmt.Println(ExpectingConnectPktError)
		return ExpectingConnectPktError
	}	
	HandleConnect(c, fh)
	for !disconnect {
		fh, err := ReadFixedHeader(c)
		// maybe use a switch here or use et method
		if err != nil {
			fmt.Println(err)
			return err
		}	
		switch fh.PktType {
		case PUBLISH:
			err = HandlePublish(c, fh)
		case SUBSCRIBE:
			err = HandleSubscribe(c, fh)
		case DISCONNECT:
			disconnect = true
		case CONNECT:
			err = SecondConnectPktError
		default:
			fmt.Println("warning: unknonw packet")
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