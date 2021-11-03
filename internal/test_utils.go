package internal

import (
	"bytes"
	"github.com/kckeiks/gopigeon/mqttlib"
	"net"
	"time"
)

// implements net.Conn interface
type testConn struct {
	b *bytes.Buffer
}

func (tc *testConn) Read(b []byte) (int, error) {
	if tc.b == nil {
		tc.b = bytes.NewBuffer([]byte{})
	}
	return tc.b.Read(b)
}

func (tc *testConn) Write(b []byte) (int, error) {
	if tc.b == nil {
		tc.b = bytes.NewBuffer([]byte{})
	}
	return tc.b.Write(b)
}

func (*testConn) Close() error { return nil }
func (*testConn) LocalAddr() net.Addr { return nil }
func (*testConn) RemoteAddr() net.Addr { return nil }
func (*testConn) SetDeadline(t time.Time) error { return nil }
func (*testConn) SetReadDeadline(t time.Time) error { return nil }
func (*testConn) SetWriteDeadline(t time.Time) error { return nil }

func NewTestClient(data []byte) *Client {
	c := &testConn{}
	c.Write(data)
	return &Client{
		Conn: c,
	}
}

func NewPingreqReqRequest() (*mqttlib.FixedHeader, []byte) {
	return &mqttlib.FixedHeader{PktType: mqttlib.Pingreq << 4}, []byte{mqttlib.Pingreq << 4, 0}
}

func NewEncodedPingres() []byte {
	return []byte{mqttlib.Pingres << 4, 0}
}

func NewEncodedDisconnect() []byte {
	return []byte{mqttlib.Disconnect << 4, 0}
}

func NewTestConnackRequest(ackFlags, rtrnCode byte) []byte {
	pktType := byte(mqttlib.Connack << 4)
	remLength := byte(2)
	ackFlags = ackFlags & 1
	return []byte{pktType, remLength, ackFlags, rtrnCode}
}

func AddTestSubscriber(c *Client, topic string) {
	SubscriberTable = &Subscribers{Subscribers: make(map[string][]*Client)}
	SubscriberTable.Subscribers[topic] = []*Client{c}
	c.Topics = append(c.Topics, topic)
}
