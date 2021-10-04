package gopigeon

import (
	"bytes"
	// "fmt"
	"net"
	"time"
	"testing"
)

// implements net.Conn interface
type testClient struct {
	b bytes.Buffer
}

func (tc *testClient) WritePacket(b []byte) {
	tc.b.Write(b)  
}

func (tc *testClient) Read(b []byte) (int, error) { 
	return tc.b.Read(b)  
}

func (tc *testClient) Write(b []byte) (int, error) {
	// We are not interested in the *ack messages
	return len(b), nil
}

func (*testClient) Close() error { return nil }
func (*testClient) LocalAddr() net.Addr { return nil }
func (*testClient) RemoteAddr() net.Addr { return nil }
func (*testClient) SetDeadline(t time.Time) error { return nil }
func (*testClient) SetReadDeadline(t time.Time) error { return nil }
func (*testClient) SetWriteDeadline(t time.Time) error { return nil }

func TestHandleConnTwoConnects(t *testing.T) {
	client := &testClient{}
	_, cp1 := newTestEncodedConnectPkt()
	_, cp2 := newTestEncodedConnectPkt()
	client.WritePacket(cp1)
	client.WritePacket(cp2)
	err := HandleConn(client)
	if err != SecondConnectPktError {
		t.Fatalf("Expected SecondConnectPktError but got %d.", err)
	}
}