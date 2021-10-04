package gopigeon

import (
	"bytes"
	"io"
	"encoding/binary"
	"sync"
	"errors"
)

var subscribers = &Subscribers{subscribers: make(map[string][]io.ReadWriter)}

type Subscribers struct {
	mu sync.Mutex
	subscribers map[string][]io.ReadWriter // = make(map[string][]io.ReadWriter)
}

type SubscribePayload struct {
	TopicFilter string
	QoS byte
}

type SubscribePacket struct {
	PacketID uint16
	Payload []SubscribePayload
}

func DecodeSubscribePacket(b []byte) (*SubscribePacket, error) {
	buf := bytes.NewBuffer(b)
	packetId := make([]byte, PACKETID_LEN)
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

func HandleSubscribe(rw io.ReadWriter, fh *FixedHeader) error {
	b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(rw, b)
    if (err != nil) {
        return err
    }
	sp, err := DecodeSubscribePacket(b)
	for _, payload := range sp.Payload {
		subscribers.addSubscriber(rw, payload.TopicFilter)
	}
	esp := EncodeSubackPacket(sp.PacketID)
	_, err = rw.Write(esp)	
	if (err != nil) {
        return err
    }
	return nil
}

func EncodeSubackPacket(pktID uint16) []byte {
    var pktType byte = SUBACK << 4
	var remLength byte = 2
	fixedHeader := []byte{pktType, remLength}
	pID := make([]byte, 2)
	binary.BigEndian.PutUint16(pID, pktID)
	return append(fixedHeader, pID...)
}

func (s *Subscribers) addSubscriber(rw io.ReadWriter, topic string) {
	s.mu.Lock()
	s.subscribers[topic] = append(s.subscribers[topic], rw)
	s.mu.Unlock()
}

func (s *Subscribers) getSubscribers(topic string) ([]io.ReadWriter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subs, ok := s.subscribers[topic]
	if ok {
		return subs, nil
	}
	return nil, errors.New("NA_TOPIC")
} 