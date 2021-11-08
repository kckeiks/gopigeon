package internal

import (
	"github.com/kckeiks/gopigeon/mqttlib"
	"io"
	"reflect"
	"testing"
	"time"
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
	c := newTestClient(pp[2:])
	// When: we handle a connnect packet and pass the ReadWriter
	err := handlePublish(c, fh)
	// Then: the ReadWriter will have the Connack pkt representation in bytes
	if err != nil {
		t.Fatalf("handlePublish failed with err %d", err)
	}
	// Then: we published the message to the subscriber
	msgPublished, err := io.ReadAll(subscriber.Conn)
	if err != nil {
		t.Fatalf("unexpected error")
	}
	// msgs are published asynchronously so we might not have the msg yet
	tryCount := 2
	for tryCount > 0 && len(msgPublished) == 0 {
		time.Sleep(time.Second)
		msgPublished, err  = io.ReadAll(subscriber.Conn)
		if err != nil {
			t.Fatalf("unexpected error")
		}
		tryCount--
	}
	if !reflect.DeepEqual(msgPublished, pp) {
		t.Fatalf("publish packet sent %+v but expected %+v.", msgPublished, pp)
	}
}
