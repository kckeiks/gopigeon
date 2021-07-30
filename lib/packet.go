package lib

import (
    "io"
	"fmt"
	"encoding/hex"

)

const (
    Connect = 1
    Connack = 2
)

type FixedHeader struct {
    pktType byte
    flags byte
    remLength byte
}

func GetControlPkt(r io.Reader) *ConnectPkt {
    buff := make([]byte, 1)
    _, err := io.ReadFull(r, buff)
    if err != nil {
        return nil
    }
    fh := FixedHeader{pktType: buff[0]}
	pt := fh.pktType >> 4
	fmt.Printf("Hex Dump:%s", hex.Dump([]byte{pt}))
	if pt == Connect {
        cp := &ConnectPkt{fixedHeader: fh}
		fmt.Printf("ConnectPkt: %+v\n", cp)
		return cp
    }
	return nil
}
