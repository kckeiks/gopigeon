package internal

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kckeiks/gopigeon/mqttlib"
)

func handlePublish(c *Client, fh *mqttlib.FixedHeader) error {
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(c.Conn, b)
	if err != nil {
		return err
	}
	p, err := mqttlib.DecodePublishPacket(bytes.NewBuffer(b), fh)
	if err != nil {
		return err
	}
	publish(p)
	return nil
}

func publish(p *mqttlib.PublishPacket) error {
	ep := mqttlib.EncodePublishPacket(p)
	subs, err := subscriberTable.GetSubscribers(p.Topic)
	if err != nil {
		return err
	}
	// TODO: if one write op goes wrong, it should not stop us from trying the rest
	for _, s := range subs {
		_, err = s.Conn.Write(ep)
		if err != nil {
			fmt.Printf("publishing error: %s\n", err.Error())
		}
	}
	return nil
}
