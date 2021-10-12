package gopigeon

import (
	"testing"
)

func TestHandleConnTwoConnects(t *testing.T) {
	// Given: a connection with a client
	c := &testConn{}
	// Given: we send two connect packets
	cp := &ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4, 
		userNameFlag:0, 
		psswdFlag:0, 
		willRetainFlag:0, 
		willQoSFlag:0, 
		willFlag:0, 
		cleanSession:1, 
		payload : []byte{0, 0},
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
		Topic: "testtopic",
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