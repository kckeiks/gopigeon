package gopigeon

import (
	"testing"
)

func TestHandleConnTwoConnects(t *testing.T) {
	// Given: a connection with a client
	c := &testConn{}
	// Given: we send two connect packets
	_, cp1 := newTestEncodedConnectPkt()
	_, cp2 := newTestEncodedConnectPkt()
	c.Write(cp1)
	c.Write(cp2)
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
	pp := NewTestEncodedPublishPkt()
	c.Write(pp)
	// When: we send this packet as the first one
	err := HandleConn(c)
	// Then: we get an error because it is a protocol violation
	if err != ExpectingConnectPktError {
		t.Fatalf("Expected ExpectingConnectPktError but got %d.", err)
	}
}