package lib

import (
    "errors"
    "io"
    "fmt"
)

const (
    Connect = 1
    Connack = 2
)


const (
    LSNibbleMask = 0x0F // 0000 1111
    MSNibbleMask = 0XF0 // 1111 0000
)

type FixedHeader struct {
    pktType byte
    flags byte
    remLength int
}

func GetRemLength(r io.Reader) (int, error) {
    // http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Table_2.4_Size    
    mul := 1
    val := 0
    encodedByte := make([]byte, 1)
    for {
        r.Read(encodedByte)
        val += int(encodedByte[0] & byte(127)) * mul
        mul *= 128
        if (mul > 128*128*128) {
            return 0, errors.New("Error: Invalid remaining length.")
        }
        if ((encodedByte[0] & byte(128)) == 0) {
            return val, nil
        }
    }
}

func GetControlPkt(r io.Reader) *ConnectPkt {
    buff := make([]byte, 1)
    _, err := io.ReadFull(r, buff)
    if err != nil {
        return nil
    }

    flags := buff[0] & LSNibbleMask 
    pktType := buff[0] >> 4
    rl, lenErr := GetRemLength(r)
    if lenErr != nil {
        return nil
    }
    fh := FixedHeader{pktType: pktType, flags: flags, remLength: rl}
	cp := newControlPkt(fh.pktType)
    cp.fh = fh
    fmt.Printf("ConnectPkt: %+v\n", cp)
	return nil
}

func newControlPkt(pktType byte) *ConnectPkt {
    switch pktType {
    case Connect:
        return &ConnectPkt{}
    default:
        return nil
    }
}

