package gopigeon

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecodeConnectPacketSuccess(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4,
		cleanSession:1,
		keepAlive: []byte{0, 0},
		payload: []byte{0, 0},
	}
	cp := decodeTestConnectPkt(expectedResult)
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

func TestDecodeConnectPacketInvalidProtocolName(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		protocolName:"INVALID",
		protocolLevel:4,
		cleanSession:1,
		keepAlive: []byte{0, 0},
	}
	cp := decodeTestConnectPkt(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// fmt.Println(result)
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
		protocolName:"MQTT",
		protocolLevel:8, // bad
		cleanSession:1,
		keepAlive: []byte{0, 0},
	}
	cp := decodeTestConnectPkt(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// fmt.Println(result)
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
	expectedResult := NewTestEncodedConnackPkt()
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Got encoded ConnackPacket %d but expected %d,", result, expectedResult)
	}
}

func TestHandleConnectPacketSuccess(t *testing.T) {
	// Given: a ReadWriter implementation like bytes.Buffer or net.Conn
	// Given: we can read our Connect packet from a ReadWriter 
	cp := decodeTestConnectPkt(&ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4,
		cleanSession:1,
		keepAlive: make([]byte, 2),
		payload: []byte{0, 0},
	})
	fh := &FixedHeader{
		PktType: Connect,
		Flags: 0, 
		RemLength: uint32(len(cp[2:])),
	}
	// Given: mqtt connection that has the pkt in its buffer
	conn := newTestMQTTConn(cp[2:])
	// When: we handle a connnect packet
	HandleConnect(&conn, fh)
	// Then: the connection will have the Connack pkt representation in bytes
	expectedResult := bytes.NewBuffer(NewTestEncodedConnackPkt())
	if !reflect.DeepEqual(conn.Conn, expectedResult) {
		t.Fatalf("Got encoded ConnackPacket %d but expected %d,", conn.Conn, expectedResult)
	}
}

func TestHandleConnectCreateClientID(t *testing.T) {
	// Given: connect pkt with client id of size zero
	cp := decodeTestConnectPkt(&ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4,
		cleanSession:1,
		keepAlive: make([]byte, 2),
		payload: []byte{0, 0},
	})
	fh := &FixedHeader{
		PktType: Connect,
		Flags: 0, 
		RemLength: uint32(len(cp[2:])),
	}
	// Given: a connection with the connect pkt
	conn := newTestMQTTConn(cp[2:])
	// When: we handle the connection
	HandleConnect(&conn, fh)
	// Then: we change the state of the connection
	// by assigning it a randomly generated client id
	if conn.ClientID == "" {
		t.Fatalf("did not expect MQTTConn.ClientID to be the empty string")
	} 
}

func TestHandleConnectValidClientID(t *testing.T) {
	// Given: connect pkt with client id of size zero
	cp := decodeTestConnectPkt(&ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4,
		cleanSession:1,
		keepAlive: []byte{0, 0},
		payload: []byte{0x00, 0x06, 0x74, 0x65, 0x73, 0x74, 0x69, 0x64},
	})
	fh := &FixedHeader{
		PktType: Connect,
		Flags: 0, 
		RemLength: uint32(len(cp[2:])),
	}
	// Given: a connection with the connect pkt
	conn := newTestMQTTConn(cp[2:])
	// When: we handle the connection
	err := HandleConnect(&conn, fh)
	// Then: there is no error
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	// Then: We set the state with the 
	expectedResult := "testid"
	if conn.ClientID != expectedResult {
		t.Fatalf("expected %s but instead got %s", expectedResult, conn.ClientID)
	} 
}

func TestHandleConnectInvalidClientID(t *testing.T) {
	// Given: connect pkt with client id of size zero
	cp := decodeTestConnectPkt(&ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4,
		cleanSession:1,
		keepAlive: make([]byte, 2),
		payload: []byte{0x00, 0x07, 0x74, 0x20, 0x65, 0x73, 0x74, 0x69, 0x64},
	})
	fh := &FixedHeader{
		PktType: Connect,
		Flags: 0, 
		RemLength: uint32(len(cp[2:])),
	}
	// Given: a connection with the connect pkt
	conn := newTestMQTTConn(cp[2:])
	// When: we handle the connection
	err := HandleConnect(&conn, fh)
	// Then: there is an error
	if err != InvalidClientIDError {
		t.Fatalf("expected InvalidClientIDError but got %s", err.Error())
	}
}
