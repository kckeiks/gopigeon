package lib

import "testing"

func TestDecodeConnectPacketSuccess(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	// When: we try to decoded it
	// Then: we get a connect packet struct with the right values
}

func TestEncodeConnackPacketSuccess(t *testing.T) {
	// Given: a decoded connack packet
	// When: we try to decode it
	// Then: we get a stream/slice of bytes that represent the pkt
}

func TestHandleConnectPacketSuccess(t *testing.T) {
	// Given: a ReadWriter implementation like bytes.Buffer or net.Conn
	// When: we handle a connnect packet and pass the ReadWriter 
	// Then: the ReadWriter will have the Connack pkt representation in bytes
}