package mqttlib

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	ProtocolLevel = 4
)

type ConnectPacket struct {
	ProtocolName   string
	ProtocolLevel  byte
	UserNameFlag   byte
	PsswdFlag      byte
	WillRetainFlag byte
	WillQoSFlag    byte
	WillFlag       byte
	CleanSession   byte
	KeepAlive      int
	ClientID       string
	WillTopic      string
	WillMsg        []byte
	Username       string
	Password       []byte
}

type ConnackPacket struct {
	AckFlags byte
	RtrnCode byte
}

func DecodeConnectPacket(b []byte) (*ConnectPacket, error) {
	var err error
	cp := &ConnectPacket{}
	buf := bytes.NewBuffer(b)
	cp.ProtocolName, err = ReadEncodedStr(buf)
	if err != nil {
		return nil, err
	}
	if cp.ProtocolName != ProtocolName {
		return nil, ProtocolNameError
	}
	cp.ProtocolLevel, err = buf.ReadByte()
	if err != nil || cp.ProtocolLevel != ProtocolLevel {
		return nil, ProtocolLevelError
	}
	connectFlags, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if (connectFlags & 1) != 0 {
		return nil, ConnectReserveFlagError
	}
	cp.UserNameFlag = (connectFlags >> 7) & 1
	cp.PsswdFlag = (connectFlags >> 6) & 1
	cp.WillRetainFlag = (connectFlags >> 5) & 1
	cp.WillQoSFlag = (connectFlags >> 3) & 3 // 3rd and 4th bits
	cp.WillFlag = (connectFlags >> 2) & 1
	cp.CleanSession = (connectFlags >> 1) & 1
	// read Keep Alive
	keepAliveBuf := make([]byte, 2)
	_, err = io.ReadFull(buf, keepAliveBuf)
	if err != nil {
		return nil, err
	}
	cp.KeepAlive = int(binary.BigEndian.Uint16(keepAliveBuf))
	// read client id
	cp.ClientID, err = ReadEncodedStr(buf)
	if err != nil {
		return nil, err
	}
	if cp.WillFlag == 1 {
		cp.WillTopic, err = ReadEncodedStr(buf)
		if err != nil {
			return nil, err
		}
		cp.WillMsg, err = ReadEncodedBytes(buf)
		if err != nil {
			return nil, err
		}
	}
	if cp.UserNameFlag == 1 {
		cp.Username, err = ReadEncodedStr(buf)
		if err != nil {
			return nil, err
		}
	}
	if cp.PsswdFlag == 1 {
		cp.Password, err = ReadEncodedBytes(buf)
		if err != nil {
			return nil, err
		}
	}
	return cp, nil
}

func EncodeConnackPacket(p ConnackPacket) []byte {
	var cp = []byte{p.AckFlags, p.RtrnCode}
	var pktType byte = Connack << 4
	var remLength byte = 2
	fixedHeader := []byte{pktType, remLength}
	return append(fixedHeader, cp...)
}
