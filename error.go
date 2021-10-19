package gopigeon

import (
	"errors"
)

var (
	ProtocolNameError        = errors.New("invalid protocol")
	ProtocolLevelError       = errors.New("invalid protocol level")
	SecondConnectPktError    = errors.New("second connect pkt received")
	ExpectingConnectPktError = errors.New("expecting a connect pkt")
	InvalidClientIDError     = errors.New("invalid client id")
	UniqueClientIDError      = errors.New("client id is not unique")
	UnknownTopicError        = errors.New("unknown topic")
)

var (
	ConnectFixedHdrReservedFlagError = errors.New("connect reserved flag in fixed header must be zero")
	ConnectReserveFlagError          = errors.New("connect reserved flag must be zero")
	PingreqReservedFlagError         = errors.New("pingreq reserved flag must be zero")
)
