package gopigeon

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

const (
	ProtocolLevel = 4
)

type ConnectPacket struct {
	protocolName   string
	protocolLevel  byte
	userNameFlag   byte
	psswdFlag      byte
	willRetainFlag byte
	willQoSFlag    byte
	willFlag       byte
	cleanSession   byte
	keepAlive      int
	clientID       string
	willTopic      string
	willMsg        []byte
	username       string
	password       []byte
}

type ConnackPacket struct {
	AckFlags byte
	RtrnCode byte
}

func DecodeConnectPacket(b []byte) (*ConnectPacket, error) {
	var err error
	cp := &ConnectPacket{}
	buf := bytes.NewBuffer(b)
	cp.protocolName, err = ReadEncodedStr(buf)
	if err != nil {
		return nil, err
	}
	if cp.protocolName != ProtocolName {
		return nil, ProtocolNameError
	}
	cp.protocolLevel, err = buf.ReadByte()
	if err != nil || cp.protocolLevel != ProtocolLevel {
		return nil, ProtocolLevelError
	}
	connectFlags, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if (connectFlags & 1) != 0 {
		return nil, ConnectReserveFlagError
	}
	cp.userNameFlag = (connectFlags >> 7) & 1
	cp.psswdFlag = (connectFlags >> 6) & 1
	cp.willRetainFlag = (connectFlags >> 5) & 1
	cp.willQoSFlag = (connectFlags >> 3) & 3 // 3rd and 4th bits
	cp.willFlag = (connectFlags >> 2) & 1
	cp.cleanSession = (connectFlags >> 1) & 1
	// read Keep Alive
	keepAliveBuf := make([]byte, 2)
	_, err = io.ReadFull(buf, keepAliveBuf)
	if err != nil {
		return nil, err
	}
	cp.keepAlive = int(binary.BigEndian.Uint16(keepAliveBuf))
	// read client id
	cp.clientID, err = ReadEncodedStr(buf)
	if err != nil {
		return nil, err
	}
	if cp.willFlag == 1 {
		cp.willTopic, err = ReadEncodedStr(buf)
		if err != nil {
			return nil, err
		}
		cp.willMsg, err = ReadEncodedBytes(buf)
		if err != nil {
			return nil, err
		}
	}
	if cp.userNameFlag == 1 {
		cp.username, err = ReadEncodedStr(buf)
		if err != nil {
			return nil, err
		}
	}
	if cp.psswdFlag == 1 {
		cp.password, err = ReadEncodedBytes(buf)
		if err != nil {
			return nil, err
		}
	}
	return cp, nil
}

func EncodeConnackPacket(p ConnackPacket) []byte {
	var cp = []byte{p.AckFlags, p.RtrnCode}
	var pktType byte = Connack << 4
	var remLength byte = 2
	fixedHeader := []byte{pktType, remLength}
	return append(fixedHeader, cp...)
}

func HandleConnect(c *MQTTConn, fh *FixedHeader) error {
	if fh.Flags != 0 {
		return ConnectFixedHdrReservedFlagError
	}
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(c.Conn, b)
	if err != nil {
		return err
	}
	connectPkt, err := DecodeConnectPacket(b)
	if err != nil {
		return err
	}
	if len(connectPkt.clientID) == 0 {
		connectPkt.clientID = NewClientID()
		// we try until we get a unique ID
		for !clientIDSet.isClientIDUnique(connectPkt.clientID) {
			connectPkt.clientID = NewClientID()
		}
	}
	if !IsValidClientID(connectPkt.clientID) {
		// Server guarantees id generated will be valid so client sent us invalid id
		// TODO: send connack with 2 error code
		return InvalidClientIDError
	}
	if !clientIDSet.isClientIDUnique(connectPkt.clientID) {
		// Server guarantees id generated will be unique so client sent us used id
		// TODO: send connack with 2 error code
		return UniqueClientIDError
	}
	clientIDSet.saveClientID(connectPkt.clientID)
	c.ClientID = connectPkt.clientID
	c.KeepAlive = connectPkt.keepAlive
	if connectPkt.keepAlive == 0 {
		c.KeepAlive = DefaultKeepAliveTime
	}
	encodedConnackPkt := EncodeConnackPacket(ConnackPacket{
		AckFlags: 0,
		RtrnCode: 0,
	})
	_, err = c.Conn.Write(encodedConnackPkt)
	if err != nil {
		return err
	}
	return nil
}

func HandleDisconnect(c *MQTTConn) {
	// remove connection from subscription table for all of topics that it subscribed to
	for _, topic := range c.Topics {
		subscribers.removeSubscriber(c, topic)
	}
	// remove ClientID from ID set
	clientIDSet.removeClientID(c.ClientID)
	c.Conn.Close()
}

func (c *MQTTConn) resetReadDeadline() {
	c.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.KeepAlive)))
}
