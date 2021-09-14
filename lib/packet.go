package lib

import (
    "errors"
    "io"
    "encoding/binary"
    "unicode/utf8"
    "unicode"
)

const (
    Connect = 1
    Connack = 2
)

const (
    LSNibbleMask = 0x0F // 0000 1111
    MSNibbleMask = 0XF0 // 1111 0000
)

const EncodedStrByteCount = 2

type FixedHeader struct {
    PktType byte
    flags byte
    RemLength uint32
}


func GetFixedHeaderFields(r io.Reader) (*FixedHeader, error) {
    buff := make([]byte, 1)
    _, err := io.ReadFull(r, buff)
    if err != nil {
        return nil, err
    }
    flags := buff[0] & LSNibbleMask 
    pktType := buff[0] >> 4

    remLength, err := GetRemLength(r)
    if err != nil {
        return nil, err
    }
    return &FixedHeader{PktType:pktType, flags:flags, RemLength:remLength}, nil
}

func GetRemLength(r io.Reader) (uint32, error) {
    // http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Table_2.4_Size    
    mul := uint64(1)
    val := uint64(0)
    maxMulVal := uint64(128*128*128)
    encodedByte := make([]byte, 1)
    for {
        _, err := r.Read(encodedByte)
        if (err != nil) {
            return 0, err
        }
        val += uint64(encodedByte[0] & byte(127)) * mul
        mul *= uint64(128)
        if (mul > maxMulVal) {
            return 0, errors.New("")
        }
        if ((encodedByte[0] & byte(128)) == 0) {
            return uint32(val), nil
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


func GetStringLength(r io.Reader) (uint16, error) {
    protocolNameLen := make([]byte, EncodedStrByteCount)
    _, err := io.ReadFull(r, protocolNameLen)
    if (err != nil) {
        return 0, errors.New("")
    }
    pnl := binary.BigEndian.Uint16(protocolNameLen)
    if (pnl != ProtocolNameLength) {
        return 0, errors.New("")
    }
    return pnl, nil
}