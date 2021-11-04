package internal

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/kckeiks/gopigeon/mqttlib"
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

func newTestClient(data []byte) *Client {
	c := &testConn{}
	c.Write(data)
	return &Client{
		Conn: c,
	}
}

func newPingreqReqRequest() (*mqttlib.FixedHeader, []byte) {
	return &mqttlib.FixedHeader{PktType: mqttlib.Pingreq << 4}, []byte{mqttlib.Pingreq << 4, 0}
}

func newEncodedPingres() []byte {
	return []byte{mqttlib.Pingres << 4, 0}
}

func newEncodedDisconnect() []byte {
	return []byte{mqttlib.Disconnect << 4, 0}
}

func newTestConnackRequest(ackFlags, rtrnCode byte) []byte {
	pktType := byte(mqttlib.Connack << 4)
	remLength := byte(2)
	ackFlags = ackFlags & 1
	return []byte{pktType, remLength, ackFlags, rtrnCode}
}

func addTestSubscriber(c *Client, topic string) {
	subscriberTable = &Subscribers{Subscribers: make(map[string][]*Client)}
	subscriberTable.Subscribers[topic] = []*Client{c}
	c.Topics = append(c.Topics, topic)
}

func newTestPublishRequest(pp mqttlib.PublishPacket) (*mqttlib.FixedHeader, []byte) {
	encodedTopic := mqttlib.EncodeStr(pp.Topic)
	// Note: When QoS is 0 there is no packet id
	remLen := uint32(len(encodedTopic) + len(pp.Payload))
	publish := []byte{mqttlib.Publish << 4}                       // flags is zero
	publish = append(publish, mqttlib.EncodeRemLength(remLen)...) // rem len
	publish = append(publish, encodedTopic...)
	// Note: When QoS is 0 there is no packet id
	publish = append(publish, pp.Payload...)
	return &mqttlib.FixedHeader{PktType: mqttlib.Publish, RemLength: remLen}, publish
}

func newTestSubscribeRequest(sp mqttlib.SubscribePacket) (*mqttlib.FixedHeader, []byte) {
	// Packet data
	packetIdBuf := make([]byte, mqttlib.PacketIDLen)
	binary.BigEndian.PutUint16(packetIdBuf, sp.PacketID)
	payloadLen := 0
	for _, p := range sp.Payload {
		// 2 is the length of the bytes that encode the str len
		// 1 is for the QoS byte
		payloadLen += 2 + len(p.TopicFilter) + 1
	}
	// header data
	pktType := mqttlib.Subscribe << 4
	pktType = pktType | 2 //reserved 0010
	remLen := uint32(len(packetIdBuf) + payloadLen)
	// build packet in bytes
	subscribe := []byte{byte(pktType)}
	subscribe = append(subscribe, mqttlib.EncodeRemLength(remLen)...)
	subscribe = append(subscribe, packetIdBuf...)
	for _, p := range sp.Payload {
		subscribe = append(subscribe, mqttlib.EncodeStr(p.TopicFilter)...)
		subscribe = append(subscribe, p.QoS)
	}
	return &mqttlib.FixedHeader{PktType: mqttlib.Subscribe, RemLength: remLen}, subscribe
}

func newTestConnectRequest(cp *mqttlib.ConnectPacket) (*mqttlib.FixedHeader, []byte) {
	connectInBytes := encodeTestConnectPkt(cp)
	// TODO: encode flags
	return &mqttlib.FixedHeader{PktType: mqttlib.Connect, RemLength: uint32(len(connectInBytes[2:]))}, connectInBytes
}

func encodeTestConnectPkt(cp *mqttlib.ConnectPacket) []byte {
	// protocol name
	pn := []byte(cp.ProtocolName)
	var pnLen = [2]byte{}
	binary.BigEndian.PutUint16(pnLen[:], uint16(len(pn)))
	// data in payload
	payload := append(make([]byte, 0), mqttlib.EncodeStr(cp.ClientID)...)
	if cp.WillTopic != "" {
		payload = append(payload, mqttlib.EncodeStr(cp.WillTopic)...)
	}
	if len(cp.WillMsg) > 0 {
		payload = append(payload, mqttlib.EncodeBytes(cp.WillMsg)...)
	}
	if cp.Username != "" {
		payload = append(payload, mqttlib.EncodeStr(cp.Username)...)
	}
	if len(cp.Password) > 0 {
		payload = append(payload, mqttlib.EncodeBytes(cp.Password)...)
	}
	keepAliveBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(keepAliveBuf, uint16(cp.KeepAlive))
	// # of bytes: 4 bytes (protocol level, connect flags, keep alive) +
	// 2 to encode protocol name len + len of bytes in protocol name + len of payload
	lenOfPkt := uint32(4 + 2 + len(pn) + len(payload))
	connect := []byte{mqttlib.Connect << 4}                         // MS nibble has type and LS nibble is reserved e.g. 0
	connect = append(connect, mqttlib.EncodeRemLength(lenOfPkt)...) // rem length
	connect = append(connect, pnLen[:]...)                          // add protocol name
	connect = append(connect, pn...)                                // add protocol name
	connect = append(connect, cp.ProtocolLevel)
	connect = append(connect, 2)               // TODO: connect flags
	connect = append(connect, keepAliveBuf...) // TODO: Keep Alive
	connect = append(connect, payload...)
	return connect
}
