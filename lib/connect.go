package lib

import (
    "bytes"
    "encoding/binary"
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
    fh *FixedHeader
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

func DecodeConnectPacket(packet []byte) error {
    buff := bytes.NewBuffer(packet)
    protocolNameLenBuff := make([]byte, EncodedStrByteCount)
    _, err := io.ReadFull(buff, protocolNameLenBuff)
    if (err != nil) {
        return errors.New("")
    }
    pnl := binary.BigEndian.Uint16(protocolNameLenBuff)
    if (pnl != ProtocolNameLength) {
        return errors.New("")
    }
    protocol := make([]byte, pnl)
    _, err = io.ReadFull(buff, protocol)
    if (err != nil) {
        return errors.New("")
    }
    if (!IsValidUTF8Encoded(protocol)) {
        return errors.New("")
    }
    var protocolLevel byte
    protocolLevel, err = buff.ReadByte()
    if (err != nil || protocolLevel != ProtocolLevel) {
        return errors.New("")
    }
    var connectFlags byte
    connectFlags, err = buff.ReadByte()
    if (err != nil) {
        return errors.New("")
    }
    keepAlive := make([]byte, KeepAliveFieldLength)
    _, err = io.ReadFull(buff, keepAlive)
    if (err != nil) {
        return errors.New("")
    }
    payload := make([]byte, len(buff.Bytes()))
    _, err = io.ReadFull(buff, payload)
    if (err != nil) {
        return errors.New("")
    }
    fmt.Println(hex.Dump(payload))
    fmt.Printf("% 08b\n", connectFlags)
    fmt.Println("Decoded.")
    return nil
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