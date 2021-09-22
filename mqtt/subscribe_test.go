package mqtt

import (
	"bytes"
	"io"
	"testing"
	"reflect"
	"encoding/binary"
)

func addTestSubscriber(topic string) io.ReadWriter {
	sub := bytes.NewBuffer([]byte{})
	subscribers = &Subscribers{subscribers: make(map[string][]io.ReadWriter)}
	subscribers.subscribers[topic] = []io.ReadWriter{sub}
	return sub
}

func TestDecodeSubscribePacketSuccess(t *testing.T) {
	// Given: a slice/stream of bytes that represent a subscribe pkt
	sp := NewTestEncodedSubscribePkt()
	// When: we decoded it
	result, err := DecodeSubscribePacket(sp[2:])
	// Then: we get a SubscribePacket struct
	expectedResult := SubscribePacket{
		PacketID: 1,
		Payload: []SubscribePayload{SubscribePayload{TopicFilter: "testtopic", QoS: 0}},
	}
	if err != nil {
		t.Fatalf("Got an unexpected error.")
	}
	if !reflect.DeepEqual(result, &expectedResult) {
		t.Fatalf("Got PublishPacket %+v but expected %+v,", result, expectedResult)
	}  
}

func TestHandleSubscribeSuccess(t *testing.T) {
	// Given: a connection/ReaderWriter with which we will be able to read a subscribe package
	sp := NewTestEncodedSubscribePkt()
	fh := &FixedHeader{
		PktType: SUBSCRIBE,
		Flags: 2, 
		RemLength: 14,
	}
	// without header
	b := bytes.NewBuffer(sp[2:]) 
	// When: we handle the connection
	err := HandleSubscribe(b, fh)
	// Then: there is no error so we assume things are ok
	if err != nil {
		t.Fatalf("HandleSubscribe failed with err %d", err)
	}
}

func TestEncodeSubackPacket(t *testing.T) {
	// Given: a packet ID
	var packetID uint16 = 2
	// When: we create an encoded suback pkt
	encodedResult := EncodeSubackPacket(packetID)
	result := binary.BigEndian.Uint16(encodedResult[2:])
	// Then: it has out packet ID
	if !reflect.DeepEqual(result, packetID) {
		t.Fatalf("EncodeSubackPacket returned %d but expected %d.", result, packetID)
	}
}

func TestAddSubscriberSuccess(t *testing.T) {
	// Given: we have a table of subscribers
	subscribers = &Subscribers{subscribers: make(map[string][]io.ReadWriter)}
	// Given: a subscriber
	sub := bytes.NewBuffer([]byte{})
	// When: we add a subscriber given a topic
	subscribers.addSubscriber(sub, "testtopic")
	// Then: we have that subscriber added to the table
	if !reflect.DeepEqual(subscribers.subscribers["testtopic"], []io.ReadWriter{sub}) {
		t.Fatalf("Subscribers, for testtopic, has %d but expected %d.", subscribers.subscribers["testtopic"], []io.ReadWriter{sub})
	}
}

func TestGetSubscribersSuccess(t *testing.T) {
	// Given: we have to table of subscriber
	// Given: an existing topic
	sub := addTestSubscriber("testtopic")
	// When: we try to find the subs for the topic in the table
	result, err := subscribers.getSubscribers("testtopic")
	// Then: We get the right list of subs
	if err != nil {
		t.Fatalf("getSubscribers returned an unexpected error %+v", err)
	}
	if !reflect.DeepEqual(result, []io.ReadWriter{sub}) {
		t.Fatalf("Subscribers, for testtopic, has %d but expected %d.", result, []io.ReadWriter{sub})
	}
}

func TestGetSubscribersFail(t *testing.T) {
	// Given: we have to table of subscriber
	addTestSubscriber("testtopic")
	// Given: a non existing topic
	unknownTopic := "foo"
	// When: we try to find the subs for the topic in the table
	_, err := subscribers.getSubscribers(unknownTopic)
	// Then: We get an error
	if err == nil {
		t.Fatalf("getSubscribers did not returned an error.")
	}
	_, ok := subscribers.subscribers[unknownTopic]
	if ok {
		t.Fatalf("Expected that there were no records for the unknown topic.")
	}
}