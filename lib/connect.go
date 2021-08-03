package lib

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

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
    fmt.Println(binary.BigEndian.Uint16(protocolNameLen))
    if (binary.BigEndian.Uint16(protocolNameLen) != 4) {
        return nil
    }
    return nil
}