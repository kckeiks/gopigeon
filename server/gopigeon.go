package server

import (
	"net"
)

func ListenAndServe() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go HandleClient(conn)

	}
}

func HandleClient(c net.Conn) {
	defer c.Close()
	var buff [512]byte
	n, readErr := c.Read(buff[0:])
	if readErr != nil {
		// handle error
	}
	_, writeErr := c.Write(buff[0:n])
	if writeErr != nil {
		// handle error
	}
}
