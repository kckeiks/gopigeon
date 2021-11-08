package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/kckeiks/gopigeon/mqttlib"
)

func handlePublish(c *Client, fh *mqttlib.FixedHeader) error {
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(bufio.NewReader(c.Conn), b)
	if err != nil {
		return err
	}
	p, err := mqttlib.DecodePublishPacket(bytes.NewBuffer(b), fh)
	if err != nil {
		return err
	}
	clients, err := subscriberTable.GetSubscribers(p.Topic)
	if err != nil {
		return err
	}
	// TODO: if one write op goes wrong, it should not stop us from trying the rest
	// TODO: send puback or pubrec dependeing on QoS
	// TODO: this loop below should be handled in a goroutine so that it doesnt block us from
	// TODO: reading responses for this client (connection)
	for _, c := range clients {
		qos := p.QoS
		s := c.findSubscription(p.Topic)
		if s.QoS < p.QoS {
			// downgrade QoS level per Subscription
			qos = s.QoS
		}
		m := mqttlib.EncodePublishPacket(&mqttlib.PublishPacket{
			QoS:      qos,
			Topic:    p.Topic,
			PacketID: p.PacketID,
			Payload:  p.Payload,
		})
		err = writeMsg(c, m)
		if err != nil {
			fmt.Printf("publishing error: %s\n", err.Error())
		}
	}
	return nil
}
