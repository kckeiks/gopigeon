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
	if tc.b == nil {
		tc.b = bytes.NewBuffer([]byte{})
	}
	return tc.b.Write(b)
}

func (*testConn) Close() error                       { return nil }
func (*testConn) LocalAddr() net.Addr                { return nil }
func (*testConn) RemoteAddr() net.Addr               { return nil }
func (*testConn) SetDeadline(t time.Time) error      { return nil }
func (*testConn) SetReadDeadline(t time.Time) error  { return nil }
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

func newPingreqRequest() []byte {
	return []byte{Pingreq << 4, 0}
}

func newEncodedPingres() []byte {
	return []byte{Pingres << 4, 0}
}

func newEncodedDisconnect() []byte {
	return []byte{Disconnect << 4, 0}
}

func newTestConnectRequest(cp *ConnectPacket) (*FixedHeader, []byte) {
	connectInBytes := encodeTestConnectPkt(cp)
	// TODO: encode flags
	return &FixedHeader{PktType: Connect, RemLength: uint32(len(connectInBytes[2:]))}, connectInBytes
}

func encodeTestConnectPkt(cp *ConnectPacket) []byte {
	// protocol name
	pn := []byte(cp.protocolName)
	var pnLen = [2]byte{}
	binary.BigEndian.PutUint16(pnLen[:], uint16(len(pn)))
	// data in payload
	payload := append(make([]byte, 0), encodeStr(cp.clientID)...)
	if cp.willTopic != "" {
		payload = append(payload, encodeStr(cp.willTopic)...)
	}
	if len(cp.willMsg) > 0 {
		payload = append(payload, encodeBytes(cp.willMsg)...)
	}
	if cp.username != "" {
		payload = append(payload, encodeStr(cp.username)...)
	}
	if len(cp.password) > 0 {
		payload = append(payload, encodeBytes(cp.password)...)
	}
	// # of bytes: 4 bytes (protocol level, connect flags, keep alive) +
	// 2 to encode protocol name len + len of bytes in protocol name + len of payload
	lenOfPkt := uint32(4 + 2 + len(pn) + len(payload))
	connect := []byte{Connect << 4}                         // MS nibble has type and LS nibble is reserved e.g. 0
	connect = append(connect, EncodeRemLength(lenOfPkt)...) // rem length
	connect = append(connect, pnLen[:]...)                  // add protocol name
	connect = append(connect, pn...)                        // add protocol name
	connect = append(connect, cp.protocolLevel)
	connect = append(connect, 2)               // TODO: connect flags
	connect = append(connect, cp.keepAlive[:]...) // TODO: Keep Alive
	connect = append(connect, payload...)
	return connect
}

func newTestConnackRequest(ackFlags, rtrnCode byte) []byte {
	pktType := byte(Connack << 4)
	remLength := byte(2)
	ackFlags = ackFlags & 1
	return []byte{pktType, remLength, ackFlags, rtrnCode}
}

func newTestPublishRequest(pp PublishPacket) (*FixedHeader, []byte) {
	encodedTopic := encodeStr(pp.Topic)
	// Note: When QoS is 0 there is no packet id
	remLen := uint32(len(encodedTopic) + len(pp.Payload))
	publish := []byte{Publish << 4}                       // flags is zero
	publish = append(publish, EncodeRemLength(remLen)...) // rem len
	publish = append(publish, encodedTopic...)
	// Note: When QoS is 0 there is no packet id
	publish = append(publish, pp.Payload...)
	return &FixedHeader{PktType: Publish, RemLength: remLen}, publish
}

func newTestSubscribeRequest(sp SubscribePacket) (*FixedHeader, []byte) {
	// Packet data
	packetIdBuf := make([]byte, PacketIDLen)
	binary.BigEndian.PutUint16(packetIdBuf, sp.PacketID)
	payloadLen := 0
	for _, p := range sp.Payload {
		// 2 is the length of the bytes that encode the str len
		// 1 is for the QoS byte
		payloadLen += 2 + len(p.TopicFilter) + 1
	}
	// header data
	pktType := Subscribe << 4
	pktType = pktType | 2 //reserved 0010
	remLen := uint32(len(packetIdBuf) + payloadLen)
	// build packet in bytes
	subscribe := []byte{byte(pktType)}
	subscribe = append(subscribe, EncodeRemLength(remLen)...)
	subscribe = append(subscribe, packetIdBuf...)
	for _, p := range sp.Payload {
		subscribe = append(subscribe, encodeStr(p.TopicFilter)...)
		subscribe = append(subscribe, p.QoS)
	}
	return &FixedHeader{PktType: Subscribe, RemLength: remLen}, subscribe
}

func addTestSubscriber(c *MQTTConn, topic string) {
	subscribers = &Subscribers{subscribers: make(map[string][]*MQTTConn)}
	subscribers.subscribers[topic] = []*MQTTConn{c}
	c.Topics = append(c.Topics, topic)
}
