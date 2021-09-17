package mqtt

import (
    "bytes"
    //"encoding/hex"
    "fmt"
    "io"
    "errors"
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
    // TODO: validate protocol name
    if (len(protocol) != PROTOCOL_NAME_LEN) {
        return nil, errors.New("")
    }
    protocolLevel, err := buf.ReadByte()
    if (err != nil || protocolLevel != PROTOCOL_LEVEL) {
        return nil, err
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
    // payload := buf.Bytes()
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
    // fmt.Println("Connect flags:")
    // fmt.Printf("% 08b\n", connectFlags)
    // fmt.Println("Payload:")
    // fmt.Println(len(payload))
    // fmt.Println(hex.Dump(payload))
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
func HandleConnectPacket(r io.ReadWriter, fh *FixedHeader) error {
    fmt.Printf("io.ReadWriter %d", r)
    fmt.Println(r)
    b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(r, b)
    if (err != nil) {
        return err
    }
    // fmt.Println("Packet without fixed header:")
    // fmt.Println(hex.Dump(b))
    cp, err := DecodeConnectPacket(b)
    if (err != nil) {
        return err
    }
    fmt.Printf("Connect Packet: %+v\n", cp)
    connackPkt := ConnackPacket{
        AckFlags: 0,
        RtrnCode: 0,
    }

    rawConnackPkt := EncodeConnackPacket(connackPkt)
    // fmt.Println("Connack Packet:")
    // fmt.Printf("%+v\n", connackPkt)
    // fmt.Println("Raw Connack Packet:")
    // fmt.Println(rawConnackPkt)
    r.Write(rawConnackPkt)
    
    return nil
}