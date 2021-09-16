package mqtt

type SubscribePackage struct {
	PacketID uint16
	TopicFilters []string
	QoS byte
}

func DecodeSubscribePackage(b []byte) (*SubscribePackage, error) {
	return nil, nil
}

func HandleSubscribe(rw io.ReadWriter) error {
	return nil
}