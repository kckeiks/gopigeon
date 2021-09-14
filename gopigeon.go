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
	for {
		fh, err := lib.GetFixedHeaderFields(c)
		// maybe use a switch here or use et method
		if err != nil {
			if err != io.EOF {
				panic(err.Error())
			}
			break
		}
		fmt.Println("Package fixed-header received: ")
		fmt.Printf("%+v\n", fh)
		switch fh.PktType {
        case lib.Connect:
            lib.HandleConnectPacket(c, fh)
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