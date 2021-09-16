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
	for {
		fh, err := lib.ReadFixedHeader(c)
		// maybe use a switch here or use et method
		if err != nil {
			if err != io.EOF {
				panic(err.Error())
			}
			break
		}
		fmt.Printf("Package fixed-header received: %+v\n", fh)
		if cs.ConnectRcvd && fh.PktType == lib.Connect {
			break
		}
		switch fh.PktType {
        case lib.Connect:
			lib.HandleConnectPacket(c, fh)
			cs.ConnectRcvd = true
		case lib.PUBLISH:
			lib.HandlePublish(c, fh)
        default:
			panic("Unknown Control Packet!")
    	}
		fmt.Println("looping...")
	}
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}