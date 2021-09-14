package lib

import (
    "bytes"
    "encoding/hex"
    "fmt"
    "io"
    "net"
)

const (   
    ProtocolName = "MQTT"
    ProtocolNameLength = 4
    ProtocolLevel = 4
    KeepAliveFieldLength = 2
    ConnackHeaderByteLength = 2
    ConnackByteLength = 4 
)

type ConnectPacket struct {
    protocolName []byte
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
    protocolNameLen, err := GetStringLength(buf)
    if err != nil {
        return nil, err
    }
    protocol := make([]byte, protocolNameLen)
    _, err = io.ReadFull(buf, protocol)
    if (err != nil) {
        return nil, err
    }
    if (!IsValidUTF8Encoded(protocol)) {
        return nil, err
    }
    protocolLevel, err := buf.ReadByte()
    if (err != nil || protocolLevel != ProtocolLevel) {
        return nil, err
    }
    connectFlags, err := buf.ReadByte()
    if (err != nil) {
        return nil, err
    }
    keepAlive := make([]byte, KeepAliveFieldLength)
    _, err = io.ReadFull(buf, keepAlive)
    if (err != nil) {
        return nil, err
    }
    payload := buf.Bytes()
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
    fmt.Println("Connect flags:")
    fmt.Printf("% 08b\n", connectFlags)
    fmt.Println("Payload:")
    fmt.Println(len(payload))
    fmt.Println(hex.Dump(payload))
    return cp, nil
}

func EncodeConnackPacket(p ConnackPacket) []byte {
    fixedHeader := []byte{byte(Connack << 4), ConnackHeaderByteLength}
    b := make([]byte, 0)
    b = append(b, fixedHeader...)
    return append(b, p.AckFlags, p.RtrnCode)
}

// TODO: Should it be other interface other than io.Reader? seems to broad
func HandleConnectPacket(r net.Conn, fh *FixedHeader) error {
    b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(r, b)
    if (err != nil) {
        return err
    }
    fmt.Println("Packet without fixed header:")
    fmt.Println(hex.Dump(b))
    p, err := DecodeConnectPacket(b)
    if (err != nil) {
        return err
    }
    fmt.Println("Control Packet:")
    fmt.Printf("%+v\n", p)

    connackPkt := ConnackPacket{
        AckFlags: 0,
        RtrnCode: 0,
    }

    rawConnackPkt := EncodeConnackPacket(connackPkt)
    fmt.Println("Connack Packet:")
    fmt.Printf("%+v\n", connackPkt)
    fmt.Println("Raw Connack Packet:")
    fmt.Println(rawConnackPkt)
    r.Write(rawConnackPkt)
    
    return nil
}