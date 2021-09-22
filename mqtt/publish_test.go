package mqtt

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDecodePublishPacket(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	cp := NewTestEncodedPublishPkt()
	// When: we try to decoded it
	// we pass the packet without the fixed header
	result, err := DecodePublishPacket(cp[2:])
	// Then: we get a connect packet struct with the right values
	expectedResult := &PublishPacket{
		Topic: "testtopic",
		Payload: []byte{116, 101, 115, 116, 109, 115, 103},
	}
	if err != nil {
		t.Fatalf("DecodePublishPacket failed with err %d", err)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Got PublishPacket %+v but expected %+v,", result, expectedResult)
	}  
}

func TestHandlePublish(t *testing.T) {
	// Given: we have a topic with a subscriber(s)
	subscriber := addTestSubscriber("testtopic")
	// Given: a ReadWriter implementation like bytes.Buffer or net.Conn
	// Given: we can read our Connect packet from a ReadWriter 
	cp := NewTestEncodedPublishPkt()
	fh := &FixedHeader{
		PktType: PUBLISH,
		Flags: 0, 
		RemLength: 18,
	}
	// without header
	rr := bytes.NewBuffer(cp[2:])
	// When: we handle a connnect packet and pass the ReadWriter
	err := HandlePublish(rr, fh)
	// Then: the ReadWriter will have the Connack pkt representation in bytes
	if err != nil {
		t.Fatalf("HandlePublish failed with err %d", err)
	}
	// Then: we published the message to the subscriber
	msgPublished, err := io.ReadAll(subscriber)
	if !reflect.DeepEqual(msgPublished, cp) {
		t.Fatalf("Publish packet sent %+v but expected %+v.", msgPublished, cp)
	}
}