package gopigeon

import (
	"errors"
	"os"
	"testing"
	"time"
)

func init() {
	
	defaultKeepAliveTime = 4
}

func TestHandleConnTwoConnects(t *testing.T) {
	// Given: a connection with a client
	c := &testConn{}
	// Given: we send two connect packets
	cp := &ConnectPacket{
		protocolName:   "MQTT",
		protocolLevel:  4,
		userNameFlag:   0,
		psswdFlag:      0,
		willRetainFlag: 0,
		willQoSFlag:    0,
		willFlag:       0,
		cleanSession:   1,
		clientID:       "",
	}
	c.Write(encodeTestConnectPkt(cp))
	c.Write(encodeTestConnectPkt(cp))
	// When: we handle this connection
	err := HandleConn(c)
	// Then: we get an error because it is a protocol violation
	if err != SecondConnectPktError {
		t.Fatalf("Expected SecondConnectPktError but got %d.", err)
	}
}

func TestHandleConnFirstPktIsNotConnect(t *testing.T) {
	// Given: a connection with a client
	c := &testConn{}
	// Given: a non-connect packet
	_, pp := newTestPublishRequest(PublishPacket{
		Topic:   "testtopic",
		Payload: []byte{0, 1},
	})
	c.Write(pp)
	// When: we send this packet as the first one
	err := HandleConn(c)
	// Then: we get an error because it is a protocol violation
	if err != ExpectingConnectPktError {
		t.Fatalf("Expected ExpectingConnectPktError but got %d.", err)
	}
}

// These tests tend to be a bit flaky if the difference is small 
// between keep alive and idle time
func TestKeepAlive(t *testing.T) {
	disconnet := newEncodedDisconnect()
	// pingreq := newPingreqRequest()
	cases := map[string]struct {
		keepAlive int
		idleTime  int
		willFail  bool
	}{
		"not going over default keep alive time": {0, 1, false},
		"not going over given keep alive time":   {5, 1, false},
		"going over default keep alive time":     {0, 12, true},
		"going over given keep alive time":       {1, 6, true},
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
