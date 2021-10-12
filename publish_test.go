package gopigeon

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDecodePublishPacket(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := PublishPacket{
		Topic: "testtopic",
		Payload: []byte{116, 101, 115, 116, 109, 115, 103},
	}
	_, pp := newTestPublishRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	result, err := DecodePublishPacket(pp[2:])
	// Then: we get a connect packet struct with the right values
	if err != nil {
		t.Fatalf("DecodePublishPacket failed with err %d", err)
	}
	if !reflect.DeepEqual(result, &expectedResult) {
		t.Fatalf("Got PublishPacket %+v but expected %+v,", result, &expectedResult)
	}  
}

func TestHandlePublish(t *testing.T) {
	// Given: we have a topic with a subscriber(s)
	subscriber := addTestSubscriber("testtopic")
	// Given: a Publish packet in bytes
	fh, pp := newTestPublishRequest(PublishPacket{
		Topic: "testtopic",
		Payload: []byte{0, 1},
	})
	// Given: we can read our Publish packet from a ReadWriter 
	rr := bytes.NewBuffer(pp[2:])  // we omit the header
	// When: we handle a connnect packet and pass the ReadWriter
	err := HandlePublish(rr, fh)
	// Then: the ReadWriter will have the Connack pkt representation in bytes
	if err != nil {
		t.Fatalf("HandlePublish failed with err %d", err)
	}
	// Then: we published the message to the subscriber
	msgPublished, err := io.ReadAll(subscriber.Conn)
	if !reflect.DeepEqual(msgPublished, pp) {
		t.Fatalf("Publish packet sent %+v but expected %+v.", msgPublished, pp)
	}
}