package internal

import (
	"github.com/kckeiks/gopigeon/mqttlib"
	"io"
)

func HandleConnect(c *Client, fh *mqttlib.FixedHeader) error {
	if fh.Flags != 0 {
		return mqttlib.ConnectFixedHdrReservedFlagError
	}
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(c.Conn, b)
	if err != nil {
		return err
	}
	connectPkt, err := mqttlib.DecodeConnectPacket(b)
	if err != nil {
		return err
	}
	if len(connectPkt.ClientID) == 0 {
		connectPkt.ClientID = NewClientID()
		// we try until we get a unique ID
		for !ClientIDSet.IsClientIDUnique(connectPkt.ClientID) {
			connectPkt.ClientID = NewClientID()
		}
	}
	if !IsValidClientID(connectPkt.ClientID) {
		// Server guarantees id generated will be valid so client sent us invalid id
		// TODO: send connack with 2 error code
		return mqttlib.InvalidClientIDError
	}
	if !ClientIDSet.IsClientIDUnique(connectPkt.ClientID) {
		// Server guarantees id generated will be unique so client sent us used id
		// TODO: send connack with 2 error code
		return mqttlib.UniqueClientIDError
	}
	ClientIDSet.SaveClientID(connectPkt.ClientID)
	c.ID = connectPkt.ClientID
	c.KeepAlive = connectPkt.KeepAlive
	if connectPkt.KeepAlive == 0 {
		c.KeepAlive = defaultKeepAliveTime
	}
	encodedConnackPkt := mqttlib.EncodeConnackPacket(mqttlib.ConnackPacket{
		AckFlags: 0,
		RtrnCode: 0,
	})
	_, err = c.Conn.Write(encodedConnackPkt)
	if err != nil {
		return err
	}
	return nil
}

func HandleDisconnect(c *Client) {
	// remove connection from subscription table for all of topics that it subscribed to
	for _, topic := range c.Topics {
		SubscriberTable.RemoveSubscriber(c, topic)
	}
	// remove ID from ID set
	ClientIDSet.RemoveClientID(c.ID)
	c.Conn.Close()
}

