package gopigeon

import (
    "bytes"
    "encoding/hex"
    "fmt"
    "io"
)

const (   
    PROTOCOL_LEVEL = 4
    KEEP_ALIVE_FIELD_LEN = 2
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
    if protocol != PROTOCOL_NAME {
        return nil, ProtocolNameError
    }
    protocolLevel, err := buf.ReadByte()
    if (err != nil || protocolLevel != PROTOCOL_LEVEL) {
        return nil, ProtocolLevelError
    }
    connectFlags, err := buf.ReadByte()
    if (err != nil) {
        return nil, err
    }
    keepAlive := make([]byte, KEEP_ALIVE_FIELD_LEN)
    _, err = io.ReadFull(buf, keepAlive)
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
    }
    return cp, nil
}

func EncodeConnackPacket(p ConnackPacket) []byte {
    var cp = []byte{p.AckFlags, p.RtrnCode}
    var pktType byte = CONNACK << 4
    var remLength byte = 2
    fixedHeader := []byte{pktType, remLength}
    return append(fixedHeader, cp...)
}

// TODO: Should it be other interface other than io.Reader? seems to broad
func HandleConnect(r io.ReadWriter, fh *FixedHeader) error {
    b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(r, b)
    if (err != nil) {
        return err
    }
    fmt.Printf("Packet without fixed header: %v\n", hex.Dump(b))
    cp, err := DecodeConnectPacket(b)
    if (err != nil) {
        return err
    }
    fmt.Printf("Decoded Connect Packet: %+v\n", cp)
    connackPkt := ConnackPacket{
        AckFlags: 0,
        RtrnCode: 0,
    }

    rawConnackPkt := EncodeConnackPacket(connackPkt)
    _, err = r.Write(rawConnackPkt)
    if (err != nil) {
        return err
    }
    return nil
}