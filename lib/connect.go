package lib

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "strings"
    "unicode/utf8"
)

const ProtocolName = "MQTT"

type ConnectPkt struct {
    pktType byte 
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

func (cp *ConnectPkt) decode(flags byte, packet []byte) error {
	buff := bytes.NewBuffer(packet)
    protocolNameLen := make([]byte, 2)
    n, err := buff.Read(protocolNameLen)
    if (n != 2 || err != nil) {
        return nil
    }
    pnl := binary.BigEndian.Uint16(protocolNameLen)
    fmt.Println(pnl)
    if (pnl != 4) {
        return nil
    }
    protocol := make([]byte, pnl)
    n, err = buff.Read(protocol)
    if (n != int(pnl) || err != nil) {
        return nil
    }
    p := string(protocol)
    // TODO: Read about strings and UTF-8
    if (!utf8.ValidString(p) || strings.Compare(p, ProtocolName) != 0) {
        return nil
    }
    fmt.Println(p)
    return nil
}