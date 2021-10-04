package gopigeon

import (
    "bytes"
    "io"
)

const (   
    ProtocolLevel = 4
    KeepAliveFieldLen = 2
)

type MQTTConnReadWritter interface {
    io.ReadWriter
}

type Connection struct {
    Conn io.ReadWriter
    Cp *ConnectPacket
    Topics []string
}

type ConnectPacket struct {
    clientID string
    protocolName string
    protocolLevel byte
    userNameFlag byte
    psswdFlag byte
    willRetainFlag byte
    willQoSFlag byte
    willFlag byte
    cleanSession byte
    keepAlive []byte    
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
    keepAlive := make([]byte, KeepAliveFieldLen)
    _, err = io.ReadFull(buf, keepAlive)
    if (err != nil) {
        return nil, err
    }
    clientID, err := ReadEncodedStr(buf)
    if err != nil {
        return nil, err
    }
    cp := &ConnectPacket{
        clientID: clientID,
        protocolName: protocol, 
        protocolLevel: protocolLevel,
        userNameFlag: connectFlags >> 7,
        psswdFlag: connectFlags >> 6,
        willRetainFlag: connectFlags >> 5,
        willQoSFlag: connectFlags >> 3,
        willFlag: connectFlags >> 2, 
        cleanSession: connectFlags >> 1,
        keepAlive: keepAlive,    
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

// TODO: Should it be other interface other than io.Reader? seems to broad
func HandleConnect(c *Connection, fh *FixedHeader) error {
    b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(c.Conn, b)
    if (err != nil) {
        return err
    }
    cp, err := DecodeConnectPacket(b)
    if (err != nil) {
        return err
    }
    c.Cp = cp
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