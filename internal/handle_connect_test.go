package internal

import (
	"github.com/kckeiks/gopigeon/mqttlib"
	"reflect"
	"testing"
)

func TestHandleConnectPacketSuccess(t *testing.T) {
	// Given: A connect request for the wire
	fh, connect := NewTestConnectRequest(&mqttlib.ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
	})
	// Given: mqtt connection that has the pkt in its buffer
	conn := NewTestClient(connect[2:])
	// When: we handle a connnect packet
	err := HandleConnect(conn, fh)
	// Then: no error
	if err != nil {
		t.Fatalf("got an unexpected error")
	}
	// Then: the connection will have the Connack pkt representation in bytes
	expectedResult := NewTestConnackRequest(0, 0)
	result := make([]byte, len(expectedResult))
	conn.Conn.Read(result)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("got encoded ConnackPacket %d but expected %d,", result, expectedResult)
	}
}

func TestHandleConnectPacketInvalidFixedHdrReservedFlag(t *testing.T) {
	// Given: A connect request for the wire
	fh, connect := NewTestConnectRequest(&mqttlib.ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
	})
	// Given: reserved flag in header is not 0
	fh.Flags = 1
	// Given: mqtt connection that has the pkt in its buffer
	conn := NewTestClient(connect[2:])
	// When: we handle a connnect packet
	err := HandleConnect(conn, fh)
	// Then: we get an error
	if err != mqttlib.ConnectFixedHdrReservedFlagError {
		t.Fatalf("expected ConnectFixedHdrReservedFlagError but instead got %s", err)
	}
}

func TestHandleConnectCreateClientID(t *testing.T) {
	// Given: encoded connect pkt with client id of size zero
	fh, connect := NewTestConnectRequest(&mqttlib.ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
	})
	// Given: a connection with the connect
	conn := NewTestClient(connect[2:])
	// When: we handle the connection
	HandleConnect(conn, fh)
	// Then: we change the state of the connection
	// by assigning it a randomly generated client id
	if conn.ID == "" {
		t.Fatalf("did not expect Client.ID to be the empty string")
	}
}

func TestHandleConnectValidClientID(t *testing.T) {
	// Given: encoded connect pkt with non-zero length client id
	decodedConnect := &mqttlib.ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
		ClientID:      "testid",
	}
	fh, connect := NewTestConnectRequest(decodedConnect)
	// Given: a connection with the connect pkt
	conn := NewTestClient(connect[2:])
	// When: we handle the connection
	err := HandleConnect(conn, fh)
	// Then: there is no error
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	// Then: We set the state with the
	expectedResult := decodedConnect.ClientID
	if conn.ID != expectedResult {
		t.Fatalf("expected %s but instead got %s", expectedResult, conn.ID)
	}
}

func TestHandleConnectInvalidClientID(t *testing.T) {
	// Given: encoded connect pkt with client id of size zero
	fh, connect := NewTestConnectRequest(&mqttlib.ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
		ClientID:      "t estid", // has space
	})
	// Given: a connection with the connect pkt
	conn := NewTestClient(connect[2:])
	// When: we handle the connection
	err := HandleConnect(conn, fh)
	// Then: there is an error
	if err != mqttlib.InvalidClientIDError {
		t.Fatalf("expected InvalidClientIDError but got %s", err.Error())
	}
}

func TestHandleConnectWhenClientIDIsNotUnique(t *testing.T) {
	// Given: client ID set with some keys
	ClientIDSet = &idSet{set: make(map[string]struct{})}
	// Given: our client ID is in the set so not unique
	clientID := "testid"
	ClientIDSet.set[clientID] = struct{}{}
	// Given: encoded connect pkt with non-zero length client id
	fh, connect := NewTestConnectRequest(&mqttlib.ConnectPacket{
		ProtocolName:  "MQTT",
		ProtocolLevel: 4,
		CleanSession:  1,
		ClientID:      clientID,
	})
	// Given: a connection with the connect pkt
	conn := NewTestClient(connect[2:])
	// When: we handle the connection
	err := HandleConnect(conn, fh)
	// Then: there is an error
	if err != mqttlib.UniqueClientIDError {
		t.Fatalf("expected UniqueClientIDError but got %s", err.Error())
	}
}

func TestHandleClose(t *testing.T) {
	// Given: we have a connection (client) in our SubscriberTable table for some topic
	topic := "sometopic"
	connection := NewTestClient([]byte{})
	connection.ID = "testclientid"
	AddTestSubscriber(connection, topic)
	// assert
	subscribed := false
	subscriberResult, _ := SubscriberTable.GetSubscribers(topic)
	for _, s := range subscriberResult {
		if s.ID == connection.ID {
			subscribed = true
		}
	}
	if !subscribed {
		t.Fatalf("given failed: client was not added to subscriber table")
	}
	// Given: we have the connection's client ID in our client ID set
	ClientIDSet = &idSet{set: make(map[string]struct{})}
	ClientIDSet.set[connection.ID] = struct{}{}
	// assert
	if ClientIDSet.IsClientIDUnique(connection.ID) {
		t.Fatalf("given failed: client id should be not be unique because it should have been added")
	}
	// When: there is a disconnect
	HandleDisconnect(connection)
	// Then: the client is not in our SubscriberTable table
	subscriberResult, err := SubscriberTable.GetSubscribers(topic)
	if err != nil {
		t.Fatalf("unexpected error %+v\n", err)
	}
	for _, s := range subscriberResult {
		if s.ID == connection.ID {
			t.Fatalf("expected connection (client) to have been removed from table")
		}
	}
	// Then: the client's client id is not in our client id set
	if !ClientIDSet.IsClientIDUnique(connection.ID) {
		t.Fatalf("client id should be unique")
	}
}
