package mqtt

import (
	"bytes"
	"io"
)

func newTestEncodedFixedHeader() io.ReadWriter {
	return bytes.NewBuffer([]byte{16, 12, 0, 4})
}

// TODO: Maybe refactor this so it looks better
func newTestEncodedConnectPkt() []byte {
	return []byte{16, 12, 0, 4, 77, 81, 84, 84, 4, 2, 0, 60, 0, 0}
}

func newTestReadWriterConnect() io.ReadWriter {
	cp := make([]byte, 0)
	cp = append(cp, newTestEncodedConnectPkt()[2:]...)
	return bytes.NewBuffer(cp)
}

func NewTestEncodedConnackPkt() []byte {
	return []byte{32, 2, 0, 0}
}

func NewTestEncodedPublishPkt() []byte {
	return []byte{48, 18, 0, 9, 116, 101, 115, 116, 116, 111, 112, 105, 99, 116, 101, 115, 116, 109, 115, 103}
}

func NewTestEncodedSubscribePkt() []byte {
	return []byte{130, 14, 0, 1, 0, 9, 116, 101, 115, 116, 116, 111, 112, 105, 99, 0}
}