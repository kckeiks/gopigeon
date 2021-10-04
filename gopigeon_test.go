package gopigeon

import (
	"bytes"
	"net"
	"time"
	"testing"
)

// implements net.Conn interface
type testConn struct {
	b bytes.Buffer
}

func (tc *testConn) TestWrite(b []byte) {
	tc.b.Write(b)  
}

func (tc *testConn) Read(b []byte) (int, error) { 
	return tc.b.Read(b)  
}

func (tc *testConn) Write(b []byte) (int, error) {
	// We are not interested in the *ack messages
	return len(b), nil
}

func (*testConn) Close() error { return nil }
func (*testConn) LocalAddr() net.Addr { return nil }
func (*testConn) RemoteAddr() net.Addr { return nil }
func (*testConn) SetDeadline(t time.Time) error { return nil }
func (*testConn) SetReadDeadline(t time.Time) error { return nil }
func (*testConn) SetWriteDeadline(t time.Time) error { return nil }

func TestHandleConnTwoConnects(t *testing.T) {
	// Given: a connection with a client
	c := &testConn{}
	// Given: we send two connect packets
	_, cp1 := newTestEncodedConnectPkt()
	_, cp2 := newTestEncodedConnectPkt()
	c.TestWrite(cp1)
	c.TestWrite(cp2)
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
	c.TestWrite(pp)
	// When: we send this packet as the first one
	err := HandleConn(c)
	// Then: we get an error because it is a protocol violation
	if err != ExpectingConnectPktError {
		t.Fatalf("Expected ExpectingConnectPktError but got %d.", err)
	}
}