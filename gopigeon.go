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
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go HandleConn(conn)
	}
}

func HandleConn(c net.Conn) error {
	connection := &MQTTConn{Conn: c}
	defer HandleDisconnect(connection)
	fh, err := ReadFixedHeader(connection.Conn)
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
		connection.resetReadDeadline()
		fh, err := ReadFixedHeader(connection.Conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch fh.PktType {
		case Publish:
			err = HandlePublish(connection.Conn, fh)
		case Subscribe:
			err = HandleSubscribe(connection, fh)
		case Pingreq:
			err = HandlePingreq(connection, fh)
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
