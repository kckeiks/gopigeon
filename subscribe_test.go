package gopigeon

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestDecodeSubscribePacketSuccess(t *testing.T) {
	// Given: a slice/stream of bytes that represent a subscribe pkt
	expectedResult := SubscribePacket{
		PacketID: 1,
		Payload:  []SubscribePayload{{TopicFilter: "testtopic", QoS: 0}},
	}
	_, sp := NewTestSubscribeRequest(expectedResult)
	// When: we decoded it
	result, err := DecodeSubscribePacket(sp[2:])
	// Then: we get a SubscribePacket struct
	if err != nil {
		t.Fatalf("Got an unexpected error.")
	}
	if !reflect.DeepEqual(result, &expectedResult) {
		t.Fatalf("Got PublishPacket %+v but expected %+v,", result, expectedResult)
	}
}

func TestHandleSubscribeSuccess(t *testing.T) {
	// Given: a connection/ReaderWriter with which we will be able to read a subscribe package
	fh, sp := NewTestSubscribeRequest(SubscribePacket{
		PacketID: 1,
		Payload:  []SubscribePayload{{TopicFilter: "testtopic", QoS: 0}},
	})
	// without header
	c := NewTestClient(sp[2:])
	// When: we handle the connection
	err := HandleSubscribe(c, fh)
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
