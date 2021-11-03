package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/kckeiks/gopigeon/mqttlib"
	"net"
	"os"
	"testing"
	"time"
)

func init() {

	defaultKeepAliveTime = 1
}

func testServer(serverError chan error) {
	ln, err := net.Listen("tcp", "localhost:8033")
	if err != nil {
		panic(fmt.Sprintf("fatal error: %v", err))
	}
	conn, err := ln.Accept()
	if err != nil {
		panic(fmt.Sprintf("fatal error: %v", err))
	}
	err = HandleClient(&Client{Conn: conn})
	ln.Close()
	serverError <- err
}

func testClient() net.Conn {
	conn, err := net.Dial("tcp", "localhost:8033")
	if err != nil {
		panic(fmt.Sprintf("failed to dial: %v\n", err))
	}
	return conn
}

func TestHandleConnTwoConnects(t *testing.T) {
	// Given: a connection with a client
	b := bytes.NewBuffer([]byte{})
	// Given: we send two connect packets
	cp := &mqttlib.ConnectPacket{
		ProtocolName:   "MQTT",
		ProtocolLevel:  4,
		UserNameFlag:   0,
		PsswdFlag:      0,
		WillRetainFlag: 0,
		WillQoSFlag:    0,
		WillFlag:       0,
		CleanSession:   1,
		ClientID:       "",
	}
	b.Write(mqttlib.EncodeTestConnectPkt(cp))
	b.Write(mqttlib.EncodeTestConnectPkt(cp))
	// When: we handle this connection
	err := HandleClient(NewTestClient(b.Bytes()))
	// Then: we get an error because it is a protocol violation
	if err != mqttlib.SecondConnectPktError {
		t.Fatalf("Expected SecondConnectPktError but got %d.", err)
	}
}

func TestHandleConnFirstPktIsNotConnect(t *testing.T) {
	// Given: a connection with a client
	b := bytes.NewBuffer([]byte{})
	// Given: a non-connect packet
	_, pp := mqttlib.NewTestPublishRequest(mqttlib.PublishPacket{
		Topic:   "testtopic",
		Payload: []byte{0, 1},
	})
	b.Write(pp)
	// When: we send this packet as the first one
	err := HandleClient(NewTestClient(b.Bytes()))
	// Then: we get an error because it is a protocol violation
	if err != mqttlib.ExpectingConnectPktError {
		t.Fatalf("Expected ExpectingConnectPktError but got %d.", err)
	}
}

// These tests tend to be a bit flaky if the difference is small
// between keep alive and idle time
func TestKeepAlive(t *testing.T) {
	disconnet := NewEncodedDisconnect()
	cases := map[string]struct {
		keepAlive int
		idleTime  int
		willFail  bool
	}{
		"not going over default keep alive time": {0, 0, false},
		"not going over given keep alive time":   {3, 1, false},
		"going over default keep alive time":     {0, 3, true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// start server
			expectedError := make(chan error, 0)
			go testServer(expectedError)
			time.Sleep(time.Second * time.Duration(1))
			// start client
			conn := testClient()
			defer conn.Close()
			_, connect := mqttlib.NewTestConnectRequest(&mqttlib.ConnectPacket{
				ProtocolName:  "MQTT",
				ProtocolLevel: 4,
				CleanSession:  1,
				KeepAlive:     tc.keepAlive,
			})
			// connect properly
			_, err := conn.Write(connect)
			if err != nil {
				t.Fatalf("unexpected error on write")
			}
			time.Sleep(time.Second * time.Duration(tc.idleTime))
			// try to disconnect
			// if setDeadline worked, server was probably blocking on Read()
			// so it didnt even get a chance to Read() disconnect
			conn.Write(disconnet)
			err = <-expectedError
			if tc.willFail && (err == nil || !errors.Is(err, os.ErrDeadlineExceeded)) {
				t.Fatalf("given %+v was expecting ErrDeadlineExceeded but got %+v", tc, err)
			} else if !tc.willFail && err != nil {
				t.Fatalf("given %+v got a unexpected error %+v\n", tc, err)
			}
		})
	}
}
