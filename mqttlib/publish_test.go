package mqttlib

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecodePublishPacket(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := PublishPacket{
		Topic:   "testtopic",
		Payload: []byte{116, 101, 115, 116, 109, 115, 103},
	}
	fh, pp := newTestPublishRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	result, err := DecodePublishPacket(bytes.NewBuffer(pp[2:]), fh)
	// Then: we get a connect packet struct with the right values
	if err != nil {
		t.Fatalf("DecodePublishPacket failed with err %d", err)
	}
	if !reflect.DeepEqual(result, &expectedResult) {
		t.Fatalf("Got PublishPacket %+v but expected %+v,", result, &expectedResult)
	}
}
