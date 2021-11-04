package mqttlib

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
	_, sp := newTestSubscribeRequest(expectedResult)
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
