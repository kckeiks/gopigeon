package gopigeon

import (
    "bytes"
    "io"
)

const (   
    ProtocolLevel = 4
    KeepAliveFieldLen = 2
)

type ConnectPacket struct {
    protocolName string
    protocolLevel byte
    userNameFlag byte
    psswdFlag byte
    willRetainFlag byte
    willQoSFlag byte
    willFlag byte
    cleanSession byte
    keepAlive [KeepAliveFieldLen]byte
    payload []byte   
}

type ConnackPacket struct {
    AckFlags byte
    RtrnCode byte
}

func DecodeConnectPacket(b []byte) (*ConnectPacket, error) {
    buf := bytes.NewBuffer(b)
    protocol, err := ReadEncodedStr(buf)
    if err != nil {
        return nil, err
    }
    if protocol != ProtocolName {
        return nil, ProtocolNameError
    }
    protocolLevel, err := buf.ReadByte()
    if (err != nil || protocolLevel != ProtocolLevel) {
        return nil, ProtocolLevelError
    }
    connectFlags, err := buf.ReadByte()
    if (err != nil) {
        return nil, err
    }
    // keepAlive := make([]byte, KeepAliveFieldLen)
    keepAlive := [KeepAliveFieldLen]byte{}
    _, err = io.ReadFull(buf, keepAlive[:])
    if (err != nil) {
        return nil, err
    }
    // read the rest
    payload := make([]byte, buf.Len())
    _, err = io.ReadFull(buf, payload)
    if (err != nil) {
        return nil, err
    }
    cp := &ConnectPacket{
        protocolName: protocol, 
        protocolLevel: protocolLevel,
        userNameFlag: connectFlags >> 7,
        psswdFlag: connectFlags >> 6,
        willRetainFlag: connectFlags >> 5,
        willQoSFlag: connectFlags >> 3,
        willFlag: connectFlags >> 2, 
        cleanSession: connectFlags >> 1,
        keepAlive: keepAlive,
        payload: payload,
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
    b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(c.Conn, b)
    if (err != nil) {
        return err
    }
    cp, err := DecodeConnectPacket(b)
    if (err != nil) {
        return err
    }
    clientID, err := ReadEncodedStr(bytes.NewBuffer(cp.payload))
    if err != nil {
        return err
    }
    if len(clientID) == 0 {
        clientID = NewClientID()
        // we try until we get a unique ID
        for !clientIDSet.isClientIDUnique(clientID) {
            clientID = NewClientID()
        }
    }
    if !IsValidClientID(clientID) {
        // Server guarantees id generated will be valid so client sent us invalid id
        // TODO: send connack with 2 error code
        return InvalidClientIDError
    }
    if !clientIDSet.isClientIDUnique(clientID) {
        // Server guarantees id generated will be unique so client sent us used id
        // TODO: send connack with 2 error code
        return UniqueClientIDError
    }
    clientIDSet.saveClientID(clientID)
    c.ClientID = clientID
    connackPkt := ConnackPacket{
        AckFlags: 0,
        RtrnCode: 0,
    }
    rawConnackPkt := EncodeConnackPacket(connackPkt)
    _, err = c.Conn.Write(rawConnackPkt)
    if (err != nil) {
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