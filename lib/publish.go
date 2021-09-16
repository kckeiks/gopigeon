package lib

import (
	"io"
)

type PublishPacket struct {
	Topic []byte
	PacketId uint16
	Payload []byte
}

func HandlePublish(rw io.ReadWriter, fh *FixedHeader) error {
	b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(rw, b)
    if (err != nil) {
        return err
    }
	return nil
}

func DecodePublishPacket(b []byte) (*PublishPacket, error) {
	return nil, nil
}


