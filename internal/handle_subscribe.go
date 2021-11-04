package internal

import (
	"io"

	"github.com/kckeiks/gopigeon/mqttlib"
)

func HandleSubscribe(c *Client, fh *mqttlib.FixedHeader) error {
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(c.Conn, b)
	if err != nil {
		return err
	}
	sp, err := mqttlib.DecodeSubscribePacket(b)
	for _, payload := range sp.Payload {
		SubscriberTable.AddSubscriber(c, payload.TopicFilter)
		c.Topics = append(c.Topics, payload.TopicFilter)
	}
	esp := mqttlib.EncodeSubackPacket(sp.PacketID)
	_, err = c.Conn.Write(esp)
	if err != nil {
		return err
	}
	return nil
}
