package lib

import (
    "bytes"
    // "encoding/binary"
    "encoding/hex"
    "errors"
    "fmt"
    "io"
)

const (   
    ProtocolName = "MQTT"
    ProtocolNameLength = 4
    ProtocolLevel = 4
    KeepAliveFieldLength = 2
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

func DecodeConnectPacket(b []byte) (*ConnectPacket, error) {
    buf := bytes.NewBuffer(b)
    protocolNameLen, err := GetStringLength(buf)
    if err != nil {
        return nil, errors.New("")
    }
    protocol := make([]byte, protocolNameLen)
    _, err = io.ReadFull(buf, protocol)
    if (err != nil) {
        return nil, errors.New("")
    }
    if (!IsValidUTF8Encoded(protocol)) {
        return nil, errors.New("")
    }
    protocolLevel, err := buf.ReadByte()
    if (err != nil || protocolLevel != ProtocolLevel) {
        return nil, errors.New("")
    }
    connectFlags, err := buf.ReadByte()
    if (err != nil) {
        return nil, errors.New("")
    }
    keepAlive := make([]byte, KeepAliveFieldLength)
    _, err = io.ReadFull(buf, keepAlive)
    if (err != nil) {
        return nil, errors.New("")
    }
    payload := make([]byte, len(buf.Bytes()))
    _, err = io.ReadFull(buf, payload)
    if (err != nil) {
        return nil, errors.New("")
    }
    fmt.Println(hex.Dump(payload))
    fmt.Printf("% 08b\n", connectFlags)
    fmt.Println("Decoded.")
    return nil, nil
}

func HandleConnectPacket(r io.Reader, fh *FixedHeader) error {
    b := make([]byte, fh.RemLength)
    _, err := io.ReadFull(r, b)
    if (err != nil) {
        return errors.New(err.Error())
    }
    DecodeConnectPacket(b)
    // handle stuff here
    return nil
}