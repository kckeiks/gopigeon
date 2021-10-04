package gopigeon

import (
	"bytes"
	// "fmt"
	"net"
	"time"
	"testing"
)

// implements net.Conn interface
type testConn struct {
	b bytes.Buffer
}

func (tc *testConn) WritePacket(b []byte) {
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
	c := &testConn{}
	_, cp1 := newTestEncodedConnectPkt()
	_, cp2 := newTestEncodedConnectPkt()
	c.WritePacket(cp1)
	c.WritePacket(cp2)
	err := HandleConn(c)
	if err != SecondConnectPktError {
		t.Fatalf("Expected SecondConnectPktError but got %d.", err)
	}
}