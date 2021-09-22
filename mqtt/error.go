package mqtt

import (
	"errors"
)

var (
	ProtocolNameError = errors.New("invalid protocol")
	ProtocolLevelError = errors.New("invalid protocol level")
)
