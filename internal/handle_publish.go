package internal

import (
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
	pp, err := mqttlib.DecodePublishPacket(b, fh)
	if err != nil {
		return err
	}
	publish(b, pp)
	return nil
}

func publish(b []byte, pp *mqttlib.PublishPacket) error {
	ep := mqttlib.EncodePublishPacket(b)
	subs, err := subscriberTable.GetSubscribers(pp.Topic)
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
