package mqttlib

import "encoding/binary"

func NewTestEncodedFixedHeader(fh FixedHeader) []byte {
	remLen := EncodeRemLength(fh.RemLength)
	pktType := fh.PktType << 4
	pktType = pktType | fh.Flags
	encodedFh := []byte{pktType}
	return append(encodedFh, remLen...)
}

func NewTestConnectRequest(cp *ConnectPacket) (*FixedHeader, []byte) {
	connectInBytes := EncodeTestConnectPkt(cp)
	// TODO: encode flags
	return &FixedHeader{PktType: Connect, RemLength: uint32(len(connectInBytes[2:]))}, connectInBytes
}

func EncodeTestConnectPkt(cp *ConnectPacket) []byte {
	// protocol name
	pn := []byte(cp.ProtocolName)
	var pnLen = [2]byte{}
	binary.BigEndian.PutUint16(pnLen[:], uint16(len(pn)))
	// data in payload
	payload := append(make([]byte, 0), EncodeStr(cp.ClientID)...)
	if cp.WillTopic != "" {
		payload = append(payload, EncodeStr(cp.WillTopic)...)
	}
	if len(cp.WillMsg) > 0 {
		payload = append(payload, EncodeBytes(cp.WillMsg)...)
	}
	if cp.Username != "" {
		payload = append(payload, EncodeStr(cp.Username)...)
	}
	if len(cp.Password) > 0 {
		payload = append(payload, EncodeBytes(cp.Password)...)
	}
	keepAliveBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(keepAliveBuf, uint16(cp.KeepAlive))
	// # of bytes: 4 bytes (protocol level, connect flags, keep alive) +
	// 2 to encode protocol name len + len of bytes in protocol name + len of payload
	lenOfPkt := uint32(4 + 2 + len(pn) + len(payload))
	connect := []byte{Connect << 4}                         // MS nibble has type and LS nibble is reserved e.g. 0
	connect = append(connect, EncodeRemLength(lenOfPkt)...) // rem length
	connect = append(connect, pnLen[:]...)                  // add protocol name
	connect = append(connect, pn...)                        // add protocol name
	connect = append(connect, cp.ProtocolLevel)
	connect = append(connect, 2)               // TODO: connect flags
	connect = append(connect, keepAliveBuf...) // TODO: Keep Alive
	connect = append(connect, payload...)
	return connect
}

func NewTestConnackRequest(ackFlags, rtrnCode byte) []byte {
	pktType := byte(Connack << 4)
	remLength := byte(2)
	ackFlags = ackFlags & 1
	return []byte{pktType, remLength, ackFlags, rtrnCode}
}

func NewTestPublishRequest(pp PublishPacket) (*FixedHeader, []byte) {
	encodedTopic := EncodeStr(pp.Topic)
	// Note: When QoS is 0 there is no packet id
	remLen := uint32(len(encodedTopic) + len(pp.Payload))
	publish := []byte{Publish << 4}                       // flags is zero
	publish = append(publish, EncodeRemLength(remLen)...) // rem len
	publish = append(publish, encodedTopic...)
	// Note: When QoS is 0 there is no packet id
	publish = append(publish, pp.Payload...)
	return &FixedHeader{PktType: Publish, RemLength: remLen}, publish
}

func NewTestSubscribeRequest(sp SubscribePacket) (*FixedHeader, []byte) {
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
		subscribe = append(subscribe, EncodeStr(p.TopicFilter)...)
		subscribe = append(subscribe, p.QoS)
	}
	return &FixedHeader{PktType: Subscribe, RemLength: remLen}, subscribe
}

