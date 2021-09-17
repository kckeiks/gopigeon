package mqtt

import (
	"bytes"
	"testing"
	"reflect"
)

func NewTestEncodedSubscribePkt() []byte {
	return []byte{130, 14, 0, 1, 0, 9, 116, 101, 115, 116, 116, 111, 112, 105, 99, 0}
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