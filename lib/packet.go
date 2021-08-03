package lib

import (
    "errors"
    "io"
    // "fmt"
    // "encoding/hex"
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

type ControlPkt interface {
    decode(flags byte, packet []byte) error
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

func GetControlPkt(r io.Reader) ControlPkt {
    buff := make([]byte, 1)

    _, err := io.ReadFull(r, buff)
    if err != nil {
        return nil
    }

    flags := buff[0] & LSNibbleMask 
    pktType := buff[0] >> 4

    rl, rlerr := GetRemLength(r)
    if rlerr != nil {
        return nil
    }

    // get the remaining bytes after fixed header
    buff = make([]byte, rl)
    _, err = io.ReadFull(r, buff)
    if err != nil {
        return nil
    }

    // fmt.Println("%s", flags)
    // fmt.Println("%s", rl)
    // fmt.Println("%s", hex.Dump(buff))
    var cp ControlPkt
    switch pktType {
    case Connect:
        cp =  &ConnectPkt{pktType: pktType}
    default:
        return nil
    }

    err = cp.decode(flags, buff)
    if err != nil {
        return nil
    }
    return cp
}
