package server

import (
	"fmt"
	"net"
	"os"
	"github.com/kckeiks/gopigeon/lib"
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

	lib.GetControlPkt(c)
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}