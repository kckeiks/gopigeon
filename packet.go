package gopigeon

import (
	"encoding/binary"
	"errors"
	"io"
	"unicode"
	"unicode/utf8"
)

const (
	Connect    = 1
	Connack    = 2
	Publish    = 3
	Subscribe  = 8
	Suback     = 9
	Disconnect = 14
)

const (
	LSNibbleMask = 0x0F // 0000 1111
	MSNibbleMask = 0xF0 // 1111 0000
)

const (
	ProtocolName    = "MQTT"
	ProtocolNameLen = 4
	PacketIDLen     = 2
	StrlenLen       = 2
	MaxClientIDLen  = 23
)

type FixedHeader struct {
	PktType   byte
	Flags     byte
	RemLength uint32
}

func ReadFixedHeader(r io.Reader) (*FixedHeader, error) {
	buff := make([]byte, 1)
	_, err := io.ReadFull(r, buff)
	if err != nil {
		return nil, err
	}
	flags := buff[0] & LSNibbleMask
	pktType := buff[0] >> 4
	remLength, err := ReadRemLength(r)
	if err != nil {
		return nil, err
	}
	return &FixedHeader{PktType: pktType, Flags: flags, RemLength: remLength}, nil
}

func ReadRemLength(r io.Reader) (uint32, error) {
	// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Table_2.4_Size
	var mul uint = 1
	var val uint = 0
	var maxMulVal uint = 128 * 128 * 128
	encodedByte := make([]byte, 1)
	for {
		_, err := r.Read(encodedByte)
		if err != nil {
			return 0, err
		}
		val += uint(encodedByte[0]&byte(127)) * mul
		if mul > maxMulVal {
			return 0, errors.New("")
		}
		mul *= 128
		if (encodedByte[0] & byte(128)) == 0 {
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
	if !utf8.Valid(bytes) {
		return false
	}
	for len(bytes) > 0 {
		r, size := utf8.DecodeRune(bytes)
		// check if it is a control character (including the null character) or if it is a non-character
		if unicode.Is(unicode.Cc, r) || unicode.Is(unicode.Noncharacter_Code_Point, r) {
			return false
		}
		bytes = bytes[:len(bytes)-size]
	}
	return true
}

func ReadEncodedBytes(r io.Reader) ([]byte, error) {
	lenBuf := make([]byte, StrlenLen)
	_, err := io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, err
	}
	b := make([]byte, binary.BigEndian.Uint16(lenBuf))
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// TODO: Fix this function
func ReadEncodedStr(r io.Reader) (string, error) {
	b := make([]byte, StrlenLen)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return "", err
	}
	str := make([]byte, binary.BigEndian.Uint16(b))
	_, err = io.ReadFull(r, str)
	if err != nil {
		return "", err
	}
	if !IsValidUTF8Encoded(str) {
		return "", errors.New("")
	}
	return string(str), nil
}

func encodeStr(s string) []byte {
	b := []byte(s)
	encodedStrLen := make([]byte, 2)
	binary.BigEndian.PutUint16(encodedStrLen, uint16(len(b))) // prefix with len
	if len(b) > int(binary.BigEndian.Uint16(encodedStrLen)) {
		panic("could not encode string")
	}
	return append(encodedStrLen, b...)
}

func encodeBytes(b []byte) []byte {
	encodedStrLen := make([]byte, 2)
	binary.BigEndian.PutUint16(encodedStrLen, uint16(len(b))) // prefix with len
	if len(b) > int(binary.BigEndian.Uint16(encodedStrLen)) {
		panic("could not encode bytes")
	}
	return append(encodedStrLen, b...)
}
