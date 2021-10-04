package gopigeon

import (
	"bytes"
	"encoding/binary"
	"io"
)

func newTestEncodedFixedHeader() io.ReadWriter {
	return bytes.NewBuffer([]byte{16, 12, 0, 4})
}

func newTestEncodedConnectPkt() (*ConnectPacket, []byte) {
	cp := &ConnectPacket{
		protocolName:"MQTT",
		protocolLevel:4, 
		userNameFlag:0, 
		psswdFlag:0, 
		willRetainFlag:0, 
		willQoSFlag:0, 
		willFlag:0, 
		cleanSession:1, 
		keepAlive:[]byte{0, 0},
	}
	return cp, decodeTestConnectPkt(cp)
}

func decodeTestConnectPkt(cp *ConnectPacket) []byte {
	// TODO: this should include header
	pn := []byte(cp.protocolName)
	var pnLen = [2]byte{}
	binary.BigEndian.PutUint16(pnLen[:], uint16(len(pn)))
	connect := []byte{16, 12}
	connect = append(connect, pnLen[:]...)  // add protocol name
	connect = append(connect, pn...)  // add protocol name
	connect = append(connect, cp.protocolLevel)
	connect = append(connect, 2) // TODO: connect flags
	connect = append(connect, 0, 0,) // TODO: Keep Alive
	return append(connect, 0, 0) // TODO: Payload
}

func NewTestEncodedConnackPkt() []byte {
	return []byte{32, 2, 0, 0}
}

// TODO: Maybe refactor these to be more configurable
func NewTestEncodedPublishPkt() []byte {
	return []byte{48, 18, 0, 9, 116, 101, 115, 116, 116, 111, 112, 105, 99, 116, 101, 115, 116, 109, 115, 103}
}

func NewTestEncodedSubscribePkt() []byte {
	return []byte{130, 14, 0, 1, 0, 9, 116, 101, 115, 116, 116, 111, 112, 105, 99, 0}
}