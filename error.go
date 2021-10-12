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
