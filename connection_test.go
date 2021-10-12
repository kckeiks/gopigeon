package gopigeon

import (
	"reflect"
	"testing"
)

func TestAddSubscriberSuccess(t *testing.T) {
	// Given: we have a table of subscribers
	subscribers = &Subscribers{subscribers: make(map[string][]*MQTTConn)}
	// Given: a subscriber
	sub := newTestMQTTConn([]byte{})
	// When: we add a subscriber given a topic
	subscribers.addSubscriber(sub, "testtopic")
	// Then: we have that subscriber added to the table
	if !reflect.DeepEqual(subscribers.subscribers["testtopic"], []*MQTTConn{sub}) {
		t.Fatalf("Subscribers, for testtopic, has %+v but expected %+v.", subscribers.subscribers["testtopic"], []*MQTTConn{sub})
	}
}

func TestGetSubscribersSuccess(t *testing.T) {
	// Given: we have to table of subscriber
	// Given: an existing topic
	sub := newTestMQTTConn([]byte{})
	addTestSubscriber(sub, "testtopic")
	// When: we try to find the subs for the topic in the table
	result, err := subscribers.getSubscribers("testtopic")
	// Then: We get the right list of subs
	if err != nil {
		t.Fatalf("getSubscribers returned an unexpected error %+v", err)
	}
	if !reflect.DeepEqual(result, []*MQTTConn{sub}) {
		t.Fatalf("Subscribers, for testtopic, has %+v but expected %+v.", result, []*MQTTConn{sub})
	}
}

func TestGetSubscribersFail(t *testing.T) {
	// Given: we have to table of subscriber
	sub := newTestMQTTConn([]byte{})
	addTestSubscriber(sub, "testtopic")
	// Given: a non existing topic
	unknownTopic := "foo"
	// When: we try to find the subs for the topic in the table
	_, err := subscribers.getSubscribers(unknownTopic)
	// Then: We get an error
	if err != UnknownTopicError {
		t.Fatalf("getSubscribers did not returned an error.")
	}
	_, ok := subscribers.subscribers[unknownTopic]
	if ok {
		t.Fatalf("Expected that there were no records for the unknown topic.")
	}
}

func TestIsValidClientIDValid(t *testing.T) {
	// Given: a sequence of valid characters for client ids
	ids := []string{"thisisvalid", "123456789098766543", "thisisvalidnitis23chars", "o"}
	// When: we ask if they're valid
	// Then: we get true
	for _, id := range ids {
		if IsValidClientID(id) != true {
			t.Fatalf("client id %s should be valid.", id)
		}
	}
}

func TestIsValidClientIDInvalid(t *testing.T) {
	// Given: a sequence of invalid characters for client ids
	ids := []string{"notvalid has space", "thisistoolonggggggggggggggggggggg",
		"thisisnotvalidbecauseithas-", "%!@#$%", " tryit ", "golang_12345", ""}
	// When: we ask if they're valid
	// Then: we get false
	for _, id := range ids {
		if IsValidClientID(id) != false {
			t.Fatalf("client id \"%s\" should be invalid.", id)
		}
	}
}

func TestWhenClientIDIsUnique(t *testing.T) {
	// Given: client ID set with some key
	clientIDSet = &idSet{set: make(map[string]struct{})}
	// Given: our client ID is not in the set
	clientID := "unique"
	// When: client ID
	if !clientIDSet.isClientIDUnique(clientID) {
		t.Fatalf("client id should be unique")
	}
}

func TestWhenClientIDIsNotUnique(t *testing.T) {
	// Given: client ID set with some key
	clientIDSet = &idSet{set: make(map[string]struct{})}
	// Given: our client ID is in the set
	clientID := "notunique"
	clientIDSet.set[clientID] = struct{}{}
	// When: client ID
	if clientIDSet.isClientIDUnique(clientID) {
		t.Fatalf("client id should not be unique")
	}
}
