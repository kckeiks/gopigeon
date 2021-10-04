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
	for {
		fh, err := ReadFixedHeader(c)
		if err != nil {
			fmt.Println(err)
			return err
		}	
		switch fh.PktType {
		case PUBLISH:
			err = HandlePublish(c, fh)
		case SUBSCRIBE:
			err = HandleSubscribe(c, fh)
		case CONNECT:
			err = SecondConnectPktError
		case DISCONNECT:
			return nil
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