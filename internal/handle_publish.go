package internal

import (
	"fmt"
	"io"

	"github.com/kckeiks/gopigeon/mqttlib"
)

func HandlePublish(rw io.ReadWriter, fh *mqttlib.FixedHeader) error {
	b := make([]byte, fh.RemLength)
	_, err := io.ReadFull(rw, b)
	if err != nil {
		return err
	}
	var hasPktID bool
	qos := (fh.Flags & 6) >> 1
	if qos == 3 {
		return mqttlib.InvalidQoSValError
	}
	if qos > 0 {
		hasPktID = true
	}
	pp, err := mqttlib.DecodePublishPacket(b, hasPktID)
	if err != nil {
		return err
	}
	ep := mqttlib.EncodePublishPacket(*fh, b)
	subs, err := SubscriberTable.GetSubscribers(pp.Topic)
	if err != nil {
		return err
	}
	// TODO: if one write op goes wrong, it should not stop us from trying the rest
	for _, s := range subs {
		_, err = s.Conn.Write(ep)
		if err != nil {
			fmt.Printf("publishing error: %s\n", err.Error())
		}
	}
	return nil
}
