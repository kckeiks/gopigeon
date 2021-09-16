package gopigeon

import (
	"fmt"
	"net"
	"os"
	"github.com/kckeiks/gopigeon/lib"
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
	cs := lib.ClientSession{}
	disconnect := false
	for !disconnect {
		fh, err := lib.ReadFixedHeader(c)
		// maybe use a switch here or use et method
		if err != nil {
			if err != io.EOF {
				panic(err.Error())
			}
			break
		}
		fmt.Printf("Package fixed-header received: %+v\n", fh)
		if cs.ConnectRcvd && fh.PktType == lib.CONNECT {
			break
		}
		switch fh.PktType {
        case lib.CONNECT:
			lib.HandleConnectPacket(c, fh)
			cs.ConnectRcvd = true
		case lib.PUBLISH:
			lib.HandlePublish(c, fh)
		case lib.DISCONNECT:
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