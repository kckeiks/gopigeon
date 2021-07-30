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
    HeaderFlagMask = 0x0F
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

    flags := buff[0] & HeaderFlagMask 
    pktType := buff[0] >> 4
    rl, lenErr := GetRemLength(r)
    if lenErr != nil {
        return nil
    }
    fh := FixedHeader{pktType: pktType, flags: flags, remLength: rl}

	if pktType == Connect {
        cp := &ConnectPkt{fixedHeader: fh}
        fmt.Printf("ConnectPkt: %+v\n", cp)
		return cp
    }
	return nil
}

