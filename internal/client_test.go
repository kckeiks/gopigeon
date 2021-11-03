package internal

import (
	"github.com/kckeiks/gopigeon/mqttlib"
	"reflect"
	"testing"
)

func TestAddSubscriberSuccess(t *testing.T) {
	// Given: we have a table of SubscriberTable
	SubscriberTable = &Subscribers{Subscribers: make(map[string][]*Client)}
	// Given: a subscriber
	sub := NewTestClient([]byte{})
	// When: we add a subscriber given a topic
	SubscriberTable.AddSubscriber(sub, "testtopic")
	// Then: we have that subscriber added to the table
	if !reflect.DeepEqual(SubscriberTable.Subscribers["testtopic"], []*Client{sub}) {
		t.Fatalf("Subscribers, for testtopic, has %+v but expected %+v.", SubscriberTable.Subscribers["testtopic"], []*Client{sub})
	}
}

func TestGetSubscribersSuccess(t *testing.T) {
	// Given: we have to table of subscriber
	// Given: an existing topic
	sub := NewTestClient([]byte{})
	AddTestSubscriber(sub, "testtopic")
	// When: we try to find the subs for the topic in the table
	result, err := SubscriberTable.GetSubscribers("testtopic")
	// Then: We get the right list of subs
	if err != nil {
		t.Fatalf("GetSubscribers returned an unexpected error %+v", err)
	}
	if !reflect.DeepEqual(result, []*Client{sub}) {
		t.Fatalf("Subscribers, for testtopic, has %+v but expected %+v.", result, []*Client{sub})
	}
}

func TestGetSubscribersFail(t *testing.T) {
	// Given: we have to table of subscriber
	sub := NewTestClient([]byte{})
	AddTestSubscriber(sub, "testtopic")
	// Given: a non existing topic
	unknownTopic := "foo"
	// When: we try to find the subs for the topic in the table
	_, err := SubscriberTable.GetSubscribers(unknownTopic)
	// Then: We get an error
	if err != mqttlib.UnknownTopicError {
		t.Fatalf("GetSubscribers did not returned an error.")
	}
	_, ok := SubscriberTable.Subscribers[unknownTopic]
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
	ClientIDSet = &idSet{set: make(map[string]struct{})}
	// Given: our client ID is not in the set
	clientID := "unique"
	// When: client ID
	if !ClientIDSet.IsClientIDUnique(clientID) {
		t.Fatalf("client id should be unique")
	}
}

func TestWhenClientIDIsNotUnique(t *testing.T) {
	// Given: client ID set with some key
	ClientIDSet = &idSet{set: make(map[string]struct{})}
	// Given: our client ID is in the set
	clientID := "notunique"
	ClientIDSet.set[clientID] = struct{}{}
	// When: client ID
	if ClientIDSet.IsClientIDUnique(clientID) {
		t.Fatalf("client id should not be unique")
	}
}
