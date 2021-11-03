package internal

import (
	"github.com/kckeiks/gopigeon/mqttlib"
	"testing"
)

func TestHandleSubscribeSuccess(t *testing.T) {
	// Given: a connection/ReaderWriter with which we will be able to read a subscribe package
	fh, sp := mqttlib.NewTestSubscribeRequest(mqttlib.SubscribePacket{
		PacketID: 1,
		Payload:  []mqttlib.SubscribePayload{{TopicFilter: "testtopic", QoS: 0}},
	})
	// without header
	c := NewTestClient(sp[2:])
	// When: we handle the connection
	err := HandleSubscribe(c, fh)
	// Then: there is no error so we assume things are ok
	if err != nil {
		t.Fatalf("HandleSubscribe failed with err %d", err)
	}
}

