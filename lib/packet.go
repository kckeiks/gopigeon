package lib

import (
    "errors"
    "io"
    "unicode/utf8"
    "unicode"
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
    remLength int // TODO: use uint64 or 32
}


type PacketDecoder interface {
    Decode(r io.Reader) error
}


func GetFixedHeaderFields(r io.Reader) (pktType byte, flags byte, remLength int, err error) {
    buff := make([]byte, 1)

    _, err = io.ReadFull(r, buff)
    if err != nil {
        return 0, 0, 0, err
    }

    flags = buff[0] & LSNibbleMask 
    pktType = buff[0] >> 4

    remLength, err = GetRemLength(r)
    if err != nil {
        return 0, 0, 0, err
    }

    return pktType, flags, remLength, nil
}

func GetRemLength(r io.Reader) (int, error) {
    // http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Table_2.4_Size    
    mul := 1
    val := 0 // TODO: use uint64 or 32
    encodedByte := make([]byte, 1)
    for {
        // TODO: add error handling here
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

func IsValidUTF8Encoded(bytes []byte) bool {
    if (!utf8.Valid(bytes)) {
        return false
    }

    for len(bytes) > 0 {
        r, size := utf8.DecodeRune(bytes)
    
        // check if it is a control character (including the null character) or if it is a non-character
        if (unicode.Is(unicode.Cc, r) || unicode.Is(unicode.Noncharacter_Code_Point, r)) {
            return false
        }
        bytes = bytes[:len(bytes)-size]
    }

    return true
}

func GetControlPkt(r io.Reader) PacketDecoder {
    pktType, flags, remLength, err := GetFixedHeaderFields(r)
    if err != nil {
        return nil
    }
    fh := &FixedHeader{pktType: pktType, flags: flags, remLength: remLength}
    var pd PacketDecoder
    switch pktType {
    case Connect:
        pd = &ConnectPkt{fh: fh}
    default:
        return nil
    }

    err = pd.Decode(r)
    if err != nil {
        return nil
    }
    return pd
}

