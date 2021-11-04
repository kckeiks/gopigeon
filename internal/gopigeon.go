package internal

import (
	"fmt"
	"net"
	"os"

	"github.com/kckeiks/gopigeon/mqttlib"
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
	fh, err := mqttlib.ReadFixedHeader(c.Conn)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if fh.PktType != mqttlib.Connect {
		fmt.Println(mqttlib.ExpectingConnectPktError)
		return mqttlib.ExpectingConnectPktError
	}
	err = HandleConnect(c, fh)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for {
		c.resetReadDeadline()
		fh, err := mqttlib.ReadFixedHeader(c.Conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch fh.PktType {
		case mqttlib.Publish:
			err = HandlePublish(c.Conn, fh)
		case mqttlib.Subscribe:
			err = HandleSubscribe(c, fh)
		case mqttlib.Pingreq:
			err = HandlePingreq(c, fh)
		case mqttlib.Connect:
			err = mqttlib.SecondConnectPktError
		case mqttlib.Disconnect:
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
