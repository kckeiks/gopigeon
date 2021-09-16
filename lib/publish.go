package lib

import (
	"io"
	"fmt"
	"bytes"
	// "encoding/binary"
	"encoding/hex"
)

type PublishPacket struct {
	Topic []byte
	PacketId uint16
	Payload []byte
}

func HandlePublish(rw io.ReadWriter, fh *FixedHeader) error {
	b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(rw, b)
    if (err != nil) {
        return err
    }
	fmt.Println("Publish Packet without fixed header:")
    fmt.Println(hex.Dump(b))
	pp, err := DecodePublishPacket(b)
	if err != nil {
		return nil
	}
	fmt.Printf("Publish Packet: %+v\n", pp)
	return nil
}

func DecodePublishPacket(b []byte) (*PublishPacket, error) {
	pktLen := len(b)
	buf := bytes.NewBuffer(b)
	topicLen, err := ReadStringLength(buf)
	topic := make([]byte, topicLen)
	_, err = io.ReadFull(buf, topic)
    if (err != nil) {
        return nil, err
    }
	// Note: When QoS is 0 there is no packet id
	payloadLen := pktLen - (STRLEN_LEN + int(topicLen))
	payload := make([]byte, payloadLen)
	_, err = io.ReadFull(buf, payload)
	if err != nil {
		return nil, err
	}
	return &PublishPacket{Topic: topic, Payload: payload}, nil
}

func (p *PublishPacket) String() string {
	return fmt.Sprintf("&PublishPacket{Topic: %s, PacketID: %d, Payload: %s}\n", string(p.Topic), p.PacketId, string(p.Payload))
}
