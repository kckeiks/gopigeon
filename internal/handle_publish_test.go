package internal

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/kckeiks/gopigeon/mqttlib"
)

func TestHandlePublish(t *testing.T) {
	// Given: we have a topic with a subscriber(s)
	subscriber := newTestClient([]byte{})
	addTestSubscriber(subscriber, "testtopic")
	// Given: a Publish packet in bytes
	fh, pp := newTestPublishRequest(mqttlib.PublishPacket{
		Topic:   "testtopic",
		Payload: []byte{0, 1},
	})
	// Given: we can read our Publish packet from a ReadWriter
	rr := bytes.NewBuffer(pp[2:]) // we omit the header
	// When: we handle a connnect packet and pass the ReadWriter
	err := handlePublish(rr, fh)
	// Then: the ReadWriter will have the Connack pkt representation in bytes
	if err != nil {
		t.Fatalf("handlePublish failed with err %d", err)
	}
	// Then: we published the message to the subscriber
	msgPublished, err := io.ReadAll(subscriber.Conn)
	if !reflect.DeepEqual(msgPublished, pp) {
		t.Fatalf("Publish packet sent %+v but expected %+v.", msgPublished, pp)
	}
}
