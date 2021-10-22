package gopigeon

import (
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

// https://medium.com/nerd-for-tech/setup-and-teardown-unit-test-in-go-bd6fa1b785cd

func init() {
	DefaultKeepAliveTime = 2
}

func server(serverError chan error) {
	ln, err := net.Listen("tcp", "localhost:8033")
	if err != nil {
		panic(fmt.Sprintf("fatal error: %v", err))
	}
	defer ln.Close()
	conn, err := ln.Accept()
	if err != nil {
		panic(fmt.Sprintf("fatal error: %v", err))
	}
	serverError <- HandleConn(conn)
	return
}

func client() net.Conn {
	conn, err := net.Dial("tcp", "localhost:8033")
	if err != nil {
		panic(fmt.Sprintf("failed to dial: %v\n", err))
	}
	return conn
}

// test 1
// Given: no keep alive time
// Given: we connect
// When: we wait longer than default keep alive time, we send a pingreq
// Then: client is disconnected

// test 2
// Given: no keep alive time
// Given: we connect
// When: we wait shorter than default keep alive time, we send a pingreq
// Then: client is not disconnected because we get a pingres

// test 3
// Given: keep alive time
// Given: we connect
// When: we wait shorter than keep alive time, we send a pingreq
// Then: client is not disconnected because we get a pingres

// test 3
// Given: keep alive time
// Given: we connect
// When: we wait longer than keep alive time, we send a pingreq
// Then: client is disconnected no pingres

func TestKeepAlive(t *testing.T) {
	disconnet := newEncodedDisconnect()
	// pingreq := newPingreqRequest()
	cases := map[string]struct {
		keepAlive       int
		pingreqWaitTime int
		willFail        bool
		err             chan error
	}{
		"not going over default keep alive time": {0, 1, false, nil},
		"not going over given keep alive time":   {2, 1, false, nil},
		"going over default keep alive time":     {0, 5, true, nil},
		"going over given keep alive time":       {2, 3, true, nil},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// start server
			tc.err = make(chan error, 0)
			go server(tc.err)
			time.Sleep(time.Second * time.Duration(1))
			// start client
			conn := client()
			defer conn.Close()
			_, connect := newTestConnectRequest(&ConnectPacket{
				protocolName:  "MQTT",
				protocolLevel: 4,
				cleanSession:  1,
				keepAlive:     tc.keepAlive,
			})
			// connect properly
			_, err := conn.Write(connect)
			if err != nil {
				t.Fatalf("unexpected error on write")
			}
			time.Sleep(time.Second * time.Duration(tc.pingreqWaitTime))
			// try to disconnect
			// if setDeadline worked, server was probably blocking on Read()
			// so it didnt even get a chance to Read() disconnect
			conn.Write(disconnet)
			err = <-tc.err
			if tc.willFail && (err == nil || !errors.Is(err, os.ErrDeadlineExceeded)) {
				t.Fatalf("was expecting ErrDeadlineExceeded but got %+v", err)
			} else if !tc.willFail && err != nil {
				t.Fatalf("unexpected failure")
			}
		})
	}
}
