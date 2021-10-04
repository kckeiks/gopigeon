package gopigeon

import (
    "errors"
    "io"
    "encoding/binary"
    "unicode/utf8"
    "unicode"
    "fmt"
)

const (
    CONNECT    = 1
    CONNACK    = 2
    PUBLISH    = 3
    SUBSCRIBE  = 8
    SUBACK     = 9
    DISCONNECT = 14
)

const (
    LSNIBBLE_MASK = 0x0F // 0000 1111
    MSNIBBLE_MASK = 0XF0 // 1111 0000
)

const (
    PROTOCOL_NAME = "MQTT"
    PROTOCOL_NAME_LEN = 4
    PACKETID_LEN = 2
    STRLEN_LEN = 2
)

type FixedHeader struct {
    PktType byte
    Flags byte
    RemLength uint32
}

func ReadFixedHeader(r io.Reader) (*FixedHeader, error) {
    buff := make([]byte, 1)
    _, err := io.ReadFull(r, buff)
    fmt.Printf("First byte: %d\n", buff)
    if err != nil {
        return nil, err
    }
    flags := buff[0] & LSNIBBLE_MASK 
    pktType := buff[0] >> 4

    remLength, err := ReadRemLength(r)
    if err != nil {
        return nil, err
    }
    return &FixedHeader{PktType:pktType, Flags:flags, RemLength:remLength}, nil
}

func ReadRemLength(r io.Reader) (uint32, error) {
    // http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Table_2.4_Size    
    var mul uint = 1
    var val uint = 0
    var maxMulVal uint = 128*128*128
    encodedByte := make([]byte, 1)
    fmt.Println("EncodedBytes: ")
    for {
        _, err := r.Read(encodedByte)
        fmt.Println(encodedByte)

        if (err != nil) {
            return 0, err
        }
        val += uint(encodedByte[0] & byte(127)) * mul
        if (mul > maxMulVal) {
            return 0, errors.New("")
        }
        mul *= 128
        if ((encodedByte[0] & byte(128)) == 0) {
            return uint32(val), nil
        }
    }
}

func EncodeRemLength(n uint32) []byte {
    result := make([]byte, 0)
    for n > 0 {
        encodedByte := byte(n % 128)
        n = n / 128
        if n > 0 {
            encodedByte = encodedByte | 128
        }
        result = append(result, encodedByte)
    }
    return result
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

// TODO: Fix this function
func ReadEncodedStr(r io.Reader) (string, error) {
    b := make([]byte, STRLEN_LEN)
    _, err := io.ReadFull(r, b)
    if (err != nil) {
        return "", err
    }
    str := make([]byte, binary.BigEndian.Uint16(b))
    _, err = io.ReadFull(r, str)
    if (err != nil) {
        return "", err
    }
    if !IsValidUTF8Encoded(str) {
        return "", errors.New("")
    }
    return string(str), nil
}
