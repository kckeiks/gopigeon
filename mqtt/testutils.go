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