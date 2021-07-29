package server

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
			// maybe log
			continue
		}
		go HandleClient(conn)

	}
}

func HandleClient(c net.Conn) {
	defer c.Close()

	var buff [512]byte
	for {
		n, readErr := c.Read(buff[0:])
		if readErr != nil {
			return
		}
		_, writeErr := c.Write(buff[0:n])
		if writeErr != nil {
			return
		}
	}
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}