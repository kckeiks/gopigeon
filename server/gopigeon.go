package server

import (
	"io"
	"log"
	"net"
)

func ListenAndServe() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// conn.Write([]byte("HEllo world!"))
		// io.Copy(conn, conn)
		go func(c net.Conn) {
			c.Write([]byte("HEllo world!"))
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}