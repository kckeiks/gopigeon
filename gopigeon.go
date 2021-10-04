package gopigeon

import (
	"fmt"
	"net"
	"os"
)

func ListenAndServe() {
	ln, err := net.Listen("tcp", ":8080")
	fmt.Println("HUH0?")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("HUH1?")
			continue
		}
		go HandleConn(conn)
	}
}

func HandleConn(c net.Conn) error {
	defer c.Close()
	fh, err := ReadFixedHeader(c)
	fmt.Printf("Fixed Header: %+v\n", fh)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if fh.PktType != Connect {
		fmt.Println(ExpectingConnectPktError)
		return ExpectingConnectPktError
	}	
	HandleConnect(c, fh)
	for {
		fh, err := ReadFixedHeader(c)
		fmt.Printf("Fixed Header: %+v\n", fh)
		if err != nil {
			fmt.Println(err)
			return err
		}	
		switch fh.PktType {
		case Publish:
			err = HandlePublish(c, fh)
		case Subscribe:
			err = HandleSubscribe(c, fh)
		case Connect:
			err = SecondConnectPktError
		case Disconnect:
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
