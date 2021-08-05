package lib

import (
    "bytes"
    "encoding/binary"
    "errors"
    "io"
    "fmt"
)

const ProtocolName = "MQTT"

type ConnectPkt struct {
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

func (cp *ConnectPkt) Decode(r io.Reader) error {
    if cp.fh == nil {
        return errors.New("")
    }

    rl := cp.fh.remLength
    packet := make([]byte, rl)
    n, err := io.ReadFull(r, packet)

    // TODO: Fix this int casting
    if (n != int(rl) || err != nil) {
        return errors.New("")
    }

	buff := bytes.NewBuffer(packet)
    protocolNameLenBuff := make([]byte, 2)

    n, err = buff.Read(protocolNameLenBuff)
    if (n != 2 || err != nil) {
        return errors.New("")
    }

    pnl := binary.BigEndian.Uint16(protocolNameLenBuff)
    if (pnl != 4) {
        return errors.New("")
    }

    protocol := make([]byte, pnl)
    n, err = buff.Read(protocol)
    if (n != int(pnl) || err != nil) {
        return errors.New("")
    }

    if (!IsValidUTF8Encoded(protocol)) {
        return errors.New("")
    }
    fmt.Printf("%s\n", protocol)
    fmt.Println(IsValidUTF8Encoded(protocol))
    
    return nil
}