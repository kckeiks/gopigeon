package internal

import (
	"bytes"
	"github.com/kckeiks/gopigeon/mqttlib"
	"io"
	"reflect"
	"testing"
)

func TestHandlePublish(t *testing.T) {
	// Given: we have a topic with a subscriber(s)
	subscriber := NewTestClient([]byte{})
	AddTestSubscriber(subscriber, "testtopic")
	// Given: a Publish packet in bytes
	fh, pp := mqttlib.NewTestPublishRequest(mqttlib.PublishPacket{
		Topic:   "testtopic",
		Payload: []byte{0, 1},
	})
	// Given: we can read our Publish packet from a ReadWriter
	rr := bytes.NewBuffer(pp[2:]) // we omit the header
	// When: we handle a connnect packet and pass the ReadWriter
	err := HandlePublish(rr, fh)
	// Then: the ReadWriter will have the Connack pkt representation in bytes
	if err != nil {
		t.Fatalf("HandlePublish failed with err %d", err)
	}
	// Then: we published the message to the subscriber
	msgPublished, err := io.ReadAll(subscriber.Conn)
	if !reflect.DeepEqual(msgPublished, pp) {
		t.Fatalf("Publish packet sent %+v but expected %+v.", msgPublished, pp)
	}
}
