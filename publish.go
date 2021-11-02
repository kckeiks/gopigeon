package gopigeon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type PublishPacket struct {
	Topic    string
	PacketID uint16
	Payload  []byte
}

func HandlePublish(rw io.ReadWriter, fh *FixedHeader) error {
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(rw, b)
	if err != nil {
		return err
	}
	var hasPktID bool
	qos := (fh.Flags & 6) >> 1
	if qos == 3 {
		return InvalidQoSValError
	}
	if qos > 0 {
		hasPktID = true
	}	
	pp, err := DecodePublishPacket(b, hasPktID)
	if err != nil {
		return err
	}
	ep := EncodePublishPacket(*fh, b)
	subs, err := subscribers.getSubscribers(pp.Topic)
	if err != nil {
		return err
	}
	// TODO: if one write op goes wrong, it should not stop us from trying the rest
	for _, s := range subs {
		_, err = s.Conn.Write(ep)
		if err != nil {
			fmt.Printf("publishing error: %s\n", err.Error())
		}
	}
	return nil
}

func DecodePublishPacket(b []byte, hasPktID bool) (*PublishPacket, error) {
	p := &PublishPacket{}
	pktLen := len(b)
	buf := bytes.NewBuffer(b)
	// get topic
	topic, err := ReadEncodedStr(buf)
	if err != nil {
		return nil, err
	}
	p.Topic = topic
	// get pkt ID
	var pktIdBuf []byte
	if hasPktID {
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

func EncodePublishPacket(fh FixedHeader, p []byte) []byte {
	var pktType byte = Publish << 4
	fixedHeader := []byte{pktType}
	fixedHeader = append(fixedHeader, EncodeRemLength(fh.RemLength)...)
	return append(fixedHeader, p...)
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

