package gopigeon

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"
)

// https://medium.com/nerd-for-tech/setup-and-teardown-unit-test-in-go-bd6fa1b785cd

func server(clientCount int) {
	ln, err := net.Listen("tcp", "localhost:8033")
	if err != nil {
		panic(fmt.Sprintf("fatal error: %v", err))
	}
	for clientCount > 0 {
		conn, err := ln.Accept()
		if err != nil {
			panic(fmt.Sprintf("fatal error: %v", err))
		}
		HandleConn(conn)
		clientCount--
	}
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
	pingreq := newPingreqRequest()
	cases := map[string]struct {
		keepAlive       uint16
		pingreqWaitTime uint16
		willFail        bool
	}{
		"success: not going over default keep alive time": {0, 2, false},
		"success: not going over given keep alive time":   {0, 1, false},
		"failure: going over default keep alive time":     {0, 4, true},
		"failure: going over given keep alive time":       {0, 3, true},
	}
	// start server
	go server(len(cases))
	time.Sleep(time.Second * time.Duration(1))

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// start client
			conn := client()
			defer conn.Close()
			connectPkt := ConnectPacket{
				protocolName:  "MQTT",
				protocolLevel: 4,
				cleanSession:  1,
			}
			binary.BigEndian.PutUint16(connectPkt.keepAlive[:], tc.keepAlive)
			_, connect := newTestConnectRequest(&connectPkt)
			_, err := conn.Write(connect)
			if err != nil {
				t.Fatalf("unexpected error on write")
			}
			time.Sleep(time.Second * time.Duration(tc.pingreqWaitTime))
			_, err = conn.Write(pingreq)
			if tc.willFail && err == nil {
				t.Logf("was expecting a failure")
				t.Fail()
			} else if !tc.willFail && err != nil {
				t.Fatalf("unexpected failure")
			}
			conn.Write(disconnet)

			for {
				_, err = conn.Write(pingreq)
				if err != nil {
					return 
				}
			}
		})
	}
}
