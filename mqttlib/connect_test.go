package mqttlib

import (
	"reflect"
	"testing"
)

func TestDecodeConnectPacketSuccess(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
	}
	_, cp := NewTestConnectRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	result, err := DecodeConnectPacket(cp[2:])
	// Then: we get a connect packet struct with the right values
	if err != nil {
		t.Fatalf("DecodeConnectPacket failed with err %d", err)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Got ConnectPackage %+v but expected %+v,", result, expectedResult)
	}
}

func TestDecodeConnectPacketInvalidReservedFlag(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	// Given: connect reserved flag bit is not zero
	_, cp := NewTestConnectRequest(&ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
	})
	cp[len(cp)-5] = 3
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// Then: we get an error
	if err != ConnectReserveFlagError {
		t.Fatalf("expected ConnectReserveFlagError but instead got %d", err)
	}
}

func TestDecodeConnectPacketInvalidProtocolName(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		ProtocolName:  "INVALID",
		ProtocolLevel: 4,
		CleanSession:  1,
	}
	_, cp := NewTestConnectRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// Then: we get a connect packet struct with the right values
	if err == nil {
		t.Fatalf("DecodeConnectPacket returned nil for error.")
	}
	if err != ProtocolNameError {
		t.Fatalf("Expected %d but got %d", ProtocolNameError, err)
	}
}

func TestDecodeConnectPacketInvalidProtocolLevel(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 8, // bad
		CleanSession:  1,
	}
	_, cp := NewTestConnectRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// Then: we get a connect packet struct with the right values
	if err == nil {
		t.Fatalf("DecodeConnectPacket returned nil for error.")
	}
	if err != ProtocolLevelError {
		t.Fatalf("Expected %d but got %d", ProtocolLevelError, err)
	}
}

func TestEncodeConnackPacketSuccess(t *testing.T) {
	// Given: a decoded connack packet
	cp := ConnackPacket{
		AckFlags: 0,
		RtrnCode: 0,
	}
	// When: we try to encode it
	result := EncodeConnackPacket(cp)
	// Then: we get a stream/slice of bytes that represent the pkt
	// expectedResult containing AckFlags: 0 and RtrnCode: 0
	expectedResult := NewTestConnackRequest(cp.AckFlags, cp.RtrnCode)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Got encoded ConnackPacket %d but expected %d,", result, expectedResult)
	}
}

