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
	pktId := p.PacketID
	switch p.QoS {
	case 0:
		break
	case 1:
		// send puback
		err = writeMsg(c, mqttlib.NewEncodedPuback(pktId))
		if err != nil {
			return err
		}
	case 2:
		// send pubrec
		err = writeMsg(c, mqttlib.NewEncodedPubrec(pktId))
		if err != nil {
			return err
		}
	default:
		return mqttlib.InvalidQoSValError

	}
	go publish(p)
	return nil
}

func publish(p *mqttlib.PublishPacket) {
	clients, err := subscriberTable.GetSubscribers(p.Topic)
	// TODO: log error
	if err != nil {
		return
	}
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
}