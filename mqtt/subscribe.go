package mqtt

import (
	"fmt"
	"bytes"
	"io"
	"encoding/binary"
	"sync"
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

type SubscribePackage struct {
	PacketID uint16
	Payload []SubscribePayload
}

func DecodeSubscribePackage(b []byte) (*SubscribePackage, error) {
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
	return &SubscribePackage{PacketID: binary.BigEndian.Uint16(packetId), Payload: payload}, nil
}

func HandleSubscribe(rw io.ReadWriter, fh *FixedHeader) error {
	b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(rw, b)
    if (err != nil) {
        return err
    }
	// fmt.Println(b)
	sp, err := DecodeSubscribePackage(b)
	fmt.Printf("Subscribe Package: %+v\n", sp)
	for _, payload := range sp.Payload {
		subscribers.addSubscriber(rw, payload.TopicFilter)
	}
	fmt.Printf("Subscribers: %+v\n", subscribers)
	return nil
}

func (s *Subscribers) addSubscriber(rw io.ReadWriter, topic string) {
	s.mu.Lock()
	s.subscribers[topic] = append(s.subscribers[topic], rw)
	s.mu.Unlock()
}