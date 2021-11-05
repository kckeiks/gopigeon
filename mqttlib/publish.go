package mqttlib

import (
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

func DecodePublishPacket(r io.Reader, fh *FixedHeader) (*PublishPacket, error) {
	p := &PublishPacket{}
	// get Flags
	p.Dup = (fh.Flags & 8) >> 3
	p.QoS = (fh.Flags & 6) >> 1
	p.Retain = fh.Flags & 1
	// get topic
	topic, err := ReadEncodedStr(r)
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
		_, err := io.ReadFull(r, pktIdBuf)
		if err != nil {
			return nil, err
		}
		p.PacketID = binary.BigEndian.Uint16(pktIdBuf)
	}
	// get payload
	payloadLen := fh.RemLength - uint32(StrlenLen+len(topic)+len(pktIdBuf))
	payload := make([]byte, payloadLen)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, err
	}
	p.Payload = payload
	return p, nil
}

func EncodePublishPacket(p *PublishPacket) []byte {
	flags := (p.Dup << 3) & 8
	flags |= (p.QoS << 1) & 6
	flags |= p.Retain & 1
	fixedHeader := []byte{(Publish << 4) | flags}
	topic := EncodeStr(p.Topic)
	var pktID []byte
	if p.QoS > 0 {
		pktID = make([]byte, 2)
		binary.BigEndian.PutUint16(pktID, p.PacketID)
	}
	remLength := len(topic) + len(pktID) + len(p.Payload)
	pkt := append(fixedHeader, EncodeRemLength(uint32(remLength))...)
	pkt = append(pkt, topic...)
	pkt = append(pkt, pktID...)
	return append(pkt, p.Payload...)
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
