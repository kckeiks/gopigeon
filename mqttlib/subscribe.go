package mqttlib

import (
	"bytes"
	"encoding/binary"
	"io"
)

type SubscribePayload struct {
	TopicFilter string
	QoS         byte
}

type SubscribePacket struct {
	PacketID uint16
	Payload  []SubscribePayload
}

func DecodeSubscribePacket(b []byte) (*SubscribePacket, error) {
	buf := bytes.NewBuffer(b)
	packetId := make([]byte, PacketIDLen)
	_, err := io.ReadFull(buf, packetId)
	if err != nil {
		return nil, err
	}
	var payload []SubscribePayload
	for buf.Len() > 0 {
		filter, err := ReadEncodedStr(buf)
		if err != nil {
			return nil, err
		}
		qos, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}
		subPayload := SubscribePayload{TopicFilter: filter, QoS: qos}
		payload = append(payload, subPayload)
	}
	return &SubscribePacket{PacketID: binary.BigEndian.Uint16(packetId), Payload: payload}, nil
}

func EncodeSubackPacket(pktID uint16) []byte {
	var pktType byte = Suback << 4
	var remLength byte = 2
	fixedHeader := []byte{pktType, remLength}
	pID := make([]byte, 2)
	binary.BigEndian.PutUint16(pID, pktID)
	return append(fixedHeader, pID...)
}
