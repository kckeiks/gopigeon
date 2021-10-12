package gopigeon

import (
	"bytes"
	"encoding/binary"
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
	// We are not interested in the *ack messages
	// return len(b), nil
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

func newTestMQTTConn(data []byte) *MQTTConn {
	c := &testConn{}
	c.Write(data)
	return &MQTTConn{
		Conn: c,
	}
}

func newTestEncodedFixedHeader(fh FixedHeader) []byte {
	remLen := EncodeRemLength(fh.RemLength)
	pktType := fh.PktType << 4
	pktType = pktType | fh.Flags
	encodedFh := []byte{pktType}
	return append(encodedFh, remLen...)
}

func newTestConnectRequest(cp *ConnectPacket) (*FixedHeader, []byte) {
	// ugly side effect but will make tests cleaner and less prone to misleading errors
	if len(cp.payload) == 0 {
		cp.payload = []byte{0, 0}
	}
	connectInBytes := encodeTestConnectPkt(cp)
	// TODO: encode flags
	return &FixedHeader{PktType: Connect, RemLength: uint32(len(connectInBytes[2:]))}, connectInBytes
}

func encodeTestConnectPkt(cp *ConnectPacket) []byte {
	pn := []byte(cp.protocolName)
	var pnLen = [2]byte{}
	binary.BigEndian.PutUint16(pnLen[:], uint16(len(pn)))
	// # of bytes: 4 bytes (protocol level, connect flags, keep alive) + 
	// 2 to encode protocol name len + len of bytes in protocol name + len of payload
	lenOfPkt := uint32(4 + 2 + len(pn) + len(cp.payload))
	connect := []byte{Connect << 4}  // MS nibble has type and LS nibble is reserved e.g. 0 
	connect = append(connect, EncodeRemLength(lenOfPkt)...)  // rem length
	connect = append(connect, pnLen[:]...)  // add protocol name
	connect = append(connect, pn...)  // add protocol name
	connect = append(connect, cp.protocolLevel)
	connect = append(connect, 2) // TODO: connect flags
	connect = append(connect, 0, 0,) // TODO: Keep Alive
	return append(connect, cp.payload...) // TODO: Payload
}

func newTestConnackRequest(ackFlags, rtrnCode byte) []byte {
	pktType := byte(Connack << 4)
	remLength := byte(2)
	ackFlags = ackFlags & 1
	return []byte{pktType, remLength, ackFlags, rtrnCode}
}

func newTestPublishRequest(pp *PublishPacket) (*FixedHeader, []byte) {
	encodedTopic := encodeStr(pp.Topic)
	// Note: When QoS is 0 there is no packet id
	remLen := uint32(len(encodedTopic) + len(pp.Payload))
	publish := []byte{Publish << 4}  // flags is zero
	publish = append(publish, EncodeRemLength(remLen)...) // rem len
	publish = append(publish, encodedTopic...)
	// Note: When QoS is 0 there is no packet id
	publish = append(publish, pp.Payload...)
	return &FixedHeader{PktType: Publish, RemLength: remLen}, publish
}

func encodeStr(s string) []byte {
	b := []byte(s)
	encodedStr := make([]byte, 2)
	binary.BigEndian.PutUint16(encodedStr, uint16(len(b)))  // prefix with len
	return append(encodedStr, b...)
}

func NewTestEncodedPublishPkt() []byte {
	return []byte{48, 18, 0, 9, 116, 101, 115, 116, 116, 111, 112, 105, 99, 116, 101, 115, 116, 109, 115, 103}
}

func NewTestEncodedSubscribePkt() []byte {
	return []byte{130, 14, 0, 1, 0, 9, 116, 101, 115, 116, 116, 111, 112, 105, 99, 0}
}