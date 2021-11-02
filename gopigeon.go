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
		go HandleClient(&Client{Conn: conn})
	}
}

func HandleClient(c *Client) error {
	defer HandleDisconnect(c)
	fh, err := ReadFixedHeader(c.Conn)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if fh.PktType != Connect {
		fmt.Println(ExpectingConnectPktError)
		return ExpectingConnectPktError
	}
	err = HandleConnect(c, fh)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for {
		c.resetReadDeadline()
		fh, err := ReadFixedHeader(c.Conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch fh.PktType {
		case Publish:
			err = HandlePublish(c.Conn, fh)
		case Subscribe:
			err = HandleSubscribe(c, fh)
		case Pingreq:
			err = HandlePingreq(c, fh)
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
}
