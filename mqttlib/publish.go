package mqttlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type PublishPacket struct {
	Dup      byte
	QoS      byte
	Retain   byte
	Topic    string
	PacketID uint16
	Payload  []byte
}

func DecodePublishPacket(b []byte, fh *FixedHeader) (*PublishPacket, error) {
	p := &PublishPacket{}
	pktLen := len(b)
	buf := bytes.NewBuffer(b)
	// get Flags
	p.Dup = (fh.Flags & 8) >> 3
	p.QoS = (fh.Flags & 6) >> 1
	p.Retain = fh.Flags & 1
	// get topic
	topic, err := ReadEncodedStr(buf)
	if err != nil {
		return nil, err
	}
	p.Topic = topic
	// get pkt ID
	var pktIdBuf []byte
	if p.QoS == 3 {
		return nil, InvalidQoSValError
	}
	if p.QoS > 0 {
		_, err := io.ReadFull(buf, pktIdBuf)
		if err != nil {
			return nil, err
		}
		p.PacketID = binary.BigEndian.Uint16(pktIdBuf)
	}
	// get payload
	payloadLen := pktLen - (StrlenLen + len(topic) + len(pktIdBuf))
	payload := make([]byte, payloadLen)
	_, err = io.ReadFull(buf, payload)
	if err != nil {
		return nil, err
	}
	p.Payload = payload
	return p, nil
}

func EncodePublishPacket(b []byte) []byte {
	fixedHeader := []byte{Publish << 4}
	fixedHeader = append(fixedHeader, EncodeRemLength(uint32(len(b)))...)
	return append(fixedHeader, b...)
}

func (p *PublishPacket) String() string {
	return fmt.Sprintf("&PublishPacket{Topic: %s, PacketID: %d, Payload: %s}\n", p.Topic, p.PacketID, string(p.Payload))
}

func NewEncodedPuback(packetID uint16) []byte {
	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, packetID)
	p := []byte{Puback << 4, 2}
	return append(p, id...)
}

func NewEncodedPubrec(packetID uint16) []byte {
	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, packetID)
	p := []byte{Pubrec << 4, 2}
	return append(p, id...)
}

func NewEncodedPubrel(packetID uint16) []byte {
	var firstByte byte = (Pubrel << 4) | 2
	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, packetID)
	p := []byte{firstByte, 2}
	return append(p, id...)
}

func NewEncodedPubcomp(packetID uint16) []byte {
	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, packetID)
	p := []byte{Pubcomp << 4, 2}
	return append(p, id...)
}
