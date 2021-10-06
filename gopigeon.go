package gopigeon

import (
	"fmt"
	"net"
	"os"
)

func ListenAndServe() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
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
	connection := &MQTTConn{Conn: c}
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
	err = HandleConnect(connection, fh)
	if err != nil {
		fmt.Println(err)
		return err
	}
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
