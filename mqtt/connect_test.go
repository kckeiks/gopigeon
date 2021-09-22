package mqtt

import (
	"bytes"
	"reflect"
	"testing"
	"fmt"
)

func NewTestEncodedConnackPkt() []byte {
	return []byte{32, 2, 0, 0}
}

func TestDecodeConnectPacketSuccess(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	cp := newTestEncodedConnectPkt()
	// When: we try to decoded it
	// we pass the packet without the fixed header
	result, err := DecodeConnectPacket(cp[2:])
	fmt.Println(result)
	// Then: we get a connect packet struct with the right values
	expectedResult := &ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4, 
		userNameFlag:0, 
		psswdFlag:0, 
		willRetainFlag:0, 
		willQoSFlag:0, 
		willFlag:0, 
		cleanSession:1, 
		keepAlive:[]byte{0, 60},
	}
	if err != nil {
		t.Fatalf("DecodeConnectPacket failed with err %d", err)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Got ConnectPackage %+v but expected %+v,", result, expectedResult)
	}  
}	

func TestEncodeConnackPacketSuccess(t *testing.T) {
	// Given: a decoded connack packet
	cp := ConnackPacket{
        AckFlags: 0,
        RtrnCode: 0,
    }
	// When: we try to decode it
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
	cp := newTestEncodedConnectPkt()
	fh := &FixedHeader{
		PktType: CONNECT,
		Flags: 0, 
		RemLength: 12,
	}
	// without header
	rr := bytes.NewBuffer(cp[2:]) 
	// When: we handle a connnect packet and pass the ReadWriter
	HandleConnect(rr, fh)
	// Then: the ReadWriter will have the Connack pkt representation in bytes
	expectedResult := bytes.NewBuffer(NewTestEncodedConnackPkt())
	if !reflect.DeepEqual(rr, expectedResult) {
		t.Fatalf("Got encoded ConnackPacket %d but expected %d,", rr, expectedResult)
	}
}