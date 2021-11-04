package internal

import (
	"testing"

	"github.com/kckeiks/gopigeon/mqttlib"
)

func TestHandleSubscribeSuccess(t *testing.T) {
	// Given: a connection/ReaderWriter with which we will be able to read a subscribe package
	fh, sp := newTestSubscribeRequest(mqttlib.SubscribePacket{
		PacketID: 1,
		Payload:  []mqttlib.SubscribePayload{{TopicFilter: "testtopic", QoS: 0}},
	})
	// without header
	c := newTestClient(sp[2:])
	// When: we handle the connection
	err := HandleSubscribe(c, fh)
	// Then: there is no error so we assume things are ok
	if err != nil {
		t.Fatalf("HandleSubscribe failed with err %d", err)
	}
}
