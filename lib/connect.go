package lib

import (
    "bytes"
    "encoding/binary"
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
    protocolNameLenBuff := make([]byte, 2)

    n, err := buff.Read(protocolNameLenBuff)
    if (n != 2 || err != nil) {
        return nil
    }

    pnl := int(binary.BigEndian.Uint16(protocolNameLenBuff))
    if (pnl != 4) {
        return nil
    }

    protocol := make([]byte, pnl)
    n, err = buff.Read(protocol)
    if (n != pnl || err != nil) {
        return nil
    }
    
    return nil
}