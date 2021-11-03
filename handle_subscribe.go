package gopigeon

import "io"

func HandleSubscribe(c *Client, fh *FixedHeader) error {
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(c.Conn, b)
	if err != nil {
		return err
	}
	sp, err := DecodeSubscribePacket(b)
	for _, payload := range sp.Payload {
		SubscriberTable.AddSubscriber(c, payload.TopicFilter)
		c.Topics = append(c.Topics, payload.TopicFilter)
	}
	esp := EncodeSubackPacket(sp.PacketID)
	_, err = c.Conn.Write(esp)
	if err != nil {
		return err
	}
	return nil
}

