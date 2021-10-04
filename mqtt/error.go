package mqtt

import (
	"errors"
)

var (
	ProtocolNameError = errors.New("invalid protocol")
	ProtocolLevelError = errors.New("invalid protocol level")
	SecondConnectPktError = errors.New("second connect pkt received")
	UnknownPktError = errors.New("unknown packet")
	ExpectingConnectPktError = errors.New("expecting a connect pkt")
)
