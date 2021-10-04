package gopigeon

import (
	"io"
	"fmt"
	"bytes"
)

type PublishPacket struct {
	Topic string
	PacketID uint16
	Payload []byte
}

func HandlePublish(rw io.ReadWriter, fh *FixedHeader) error {
	b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(rw, b)
    if (err != nil) {
        return err
    }
	pp, err := DecodePublishPacket(b)
	if err != nil {
		return err
	}
	ep := EncodePublishPacket(*fh, b)
	subs, err := subscribers.getSubscribers(pp.Topic)
	if err != nil {
		return err
	}
	for _, s := range subs {
		_, err = s.Write(ep)
		if err != nil {
			return err
		}
	}
	return nil
}

func DecodePublishPacket(b []byte) (*PublishPacket, error) {
	pktLen := len(b)
	buf := bytes.NewBuffer(b)
	topic, err := ReadEncodedStr(buf)
    if (err != nil) {
        return nil, err
    }
	// Note: When QoS is 0 there is no packet id
	payloadLen := pktLen - (STRLEN_LEN + len(topic))
	payload := make([]byte, payloadLen)
	_, err = io.ReadFull(buf, payload)
	if err != nil {
		return nil, err
	}
	return &PublishPacket{Topic: topic, Payload: payload}, nil
}

func EncodePublishPacket(fh FixedHeader, p []byte) []byte {
    var pktType byte = PUBLISH << 4
    fixedHeader := []byte{pktType}
	fixedHeader = append(fixedHeader, EncodeRemLength(fh.RemLength)...)
    return append(fixedHeader, p...)
}

func (p *PublishPacket) String() string {
	return fmt.Sprintf("&PublishPacket{Topic: %s, PacketID: %d, Payload: %s}\n", p.Topic, p.PacketID, string(p.Payload))
}
