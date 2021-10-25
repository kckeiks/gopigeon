package gopigeon

import (
	"reflect"
	"testing"
)

func TestDecodeConnectPacketSuccess(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
	}
	_, cp := newTestConnectRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	result, err := DecodeConnectPacket(cp[2:])
	// Then: we get a connect packet struct with the right values
	if err != nil {
		t.Fatalf("DecodeConnectPacket failed with err %d", err)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Got ConnectPackage %+v but expected %+v,", result, expectedResult)
	}
}

func TestDecodeConnectPacketInvalidReservedFlag(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	// Given: connect reserved flag bit is not zero
	_, cp := newTestConnectRequest(&ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
	})
	cp[len(cp)-5] = 3
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// Then: we get an error
	if err != ConnectReserveFlagError {
		t.Fatalf("expected ConnectReserveFlagError but instead got %d", err)
	}
}

func TestDecodeConnectPacketInvalidProtocolName(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		protocolName:  "INVALID",
		protocolLevel: 4,
		cleanSession:  1,
	}
	_, cp := newTestConnectRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// Then: we get a connect packet struct with the right values
	if err == nil {
		t.Fatalf("DecodeConnectPacket returned nil for error.")
	}
	if err != ProtocolNameError {
		t.Fatalf("Expected %d but got %d", ProtocolNameError, err)
	}
}

func TestDecodeConnectPacketInvalidProtocolLevel(t *testing.T) {
	// Given: a stream/slice of bytes that represents a connect pkt
	expectedResult := &ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 8, // bad
		cleanSession:  1,
	}
	_, cp := newTestConnectRequest(expectedResult)
	// When: we try to decoded it
	// we pass the packet without the fixed header
	_, err := DecodeConnectPacket(cp[2:])
	// Then: we get a connect packet struct with the right values
	if err == nil {
		t.Fatalf("DecodeConnectPacket returned nil for error.")
	}
	if err != ProtocolLevelError {
		t.Fatalf("Expected %d but got %d", ProtocolLevelError, err)
	}
}

func TestEncodeConnackPacketSuccess(t *testing.T) {
	// Given: a decoded connack packet
	cp := ConnackPacket{
		AckFlags: 0,
		RtrnCode: 0,
	}
	// When: we try to encode it
	result := EncodeConnackPacket(cp)
	// Then: we get a stream/slice of bytes that represent the pkt
	// expectedResult containing AckFlags: 0 and RtrnCode: 0
	expectedResult := newTestConnackRequest(cp.AckFlags, cp.RtrnCode)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Got encoded ConnackPacket %d but expected %d,", result, expectedResult)
	}
}

func TestHandleConnectPacketSuccess(t *testing.T) {
	// Given: A connect request for the wire
	fh, connect := newTestConnectRequest(&ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
	})
	// Given: mqtt connection that has the pkt in its buffer
	conn := newTestMQTTConn(connect[2:])
	// When: we handle a connnect packet
	err := HandleConnect(conn, fh)
	// Then: no error
	if err != nil {
		t.Fatalf("got an unexpected error")
	}
	// Then: the connection will have the Connack pkt representation in bytes
	expectedResult := newTestConnackRequest(0, 0)
	result := make([]byte, len(expectedResult))
	conn.Conn.Read(result)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("got encoded ConnackPacket %d but expected %d,", result, expectedResult)
	}
}

func TestHandleConnectPacketInvalidFixedHdrReservedFlag(t *testing.T) {
	// Given: A connect request for the wire
	fh, connect := newTestConnectRequest(&ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
	})
	// Given: reserved flag in header is not 0
	fh.Flags = 1
	// Given: mqtt connection that has the pkt in its buffer
	conn := newTestMQTTConn(connect[2:])
	// When: we handle a connnect packet
	err := HandleConnect(conn, fh)
	// Then: we get an error
	if err != ConnectFixedHdrReservedFlagError {
		t.Fatalf("expected ConnectFixedHdrReservedFlagError but instead got %s", err)
	}
}

func TestHandleConnectCreateClientID(t *testing.T) {
	// Given: encoded connect pkt with client id of size zero
	fh, connect := newTestConnectRequest(&ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
	})
	// Given: a connection with the connect
	conn := newTestMQTTConn(connect[2:])
	// When: we handle the connection
	HandleConnect(conn, fh)
	// Then: we change the state of the connection
	// by assigning it a randomly generated client id
	if conn.ClientID == "" {
		t.Fatalf("did not expect MQTTConn.ClientID to be the empty string")
	}
}

func TestHandleConnectValidClientID(t *testing.T) {
	// Given: encoded connect pkt with non-zero length client id
	decodedConnect := &ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
		clientID:      "testid",
	}
	fh, connect := newTestConnectRequest(decodedConnect)
	// Given: a connection with the connect pkt
	conn := newTestMQTTConn(connect[2:])
	// When: we handle the connection
	err := HandleConnect(conn, fh)
	// Then: there is no error
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	// Then: We set the state with the
	expectedResult := decodedConnect.clientID
	if conn.ClientID != expectedResult {
		t.Fatalf("expected %s but instead got %s", expectedResult, conn.ClientID)
	}
}

func TestHandleConnectInvalidClientID(t *testing.T) {
	// Given: encoded connect pkt with client id of size zero
	fh, connect := newTestConnectRequest(&ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
		clientID:      "t estid", // has space
	})
	// Given: a connection with the connect pkt
	conn := newTestMQTTConn(connect[2:])
	// When: we handle the connection
	err := HandleConnect(conn, fh)
	// Then: there is an error
	if err != InvalidClientIDError {
		t.Fatalf("expected InvalidClientIDError but got %s", err.Error())
	}
}

func TestHandleConnectWhenClientIDIsNotUnique(t *testing.T) {
	// Given: client ID set with some keys
	clientIDSet = &idSet{set: make(map[string]struct{})}
	// Given: our client ID is in the set so not unique
	clientID := "testid"
	clientIDSet.set[clientID] = struct{}{}
	// Given: encoded connect pkt with non-zero length client id
	fh, connect := newTestConnectRequest(&ConnectPacket{
		protocolName:  "MQTT",
		protocolLevel: 4,
		cleanSession:  1,
		clientID:      clientID,
	})
	// Given: a connection with the connect pkt
	conn := newTestMQTTConn(connect[2:])
	// When: we handle the connection
	err := HandleConnect(conn, fh)
	// Then: there is an error
	if err != UniqueClientIDError {
		t.Fatalf("expected UniqueClientIDError but got %s", err.Error())
	}
}

func TestHandleClose(t *testing.T) {
	// Given: we have a connection (client) in our subscribers table for some topic
	topic := "sometopic"
	connection := newTestMQTTConn([]byte{})
	connection.ClientID = "testclientid"
	addTestSubscriber(connection, topic)
	// assert
	subscribed := false
	subscriberResult, _ := subscribers.getSubscribers(topic)
	for _, s := range subscriberResult {
		if s.ClientID == connection.ClientID {
			subscribed = true
		}
	}
	if !subscribed {
		t.Fatalf("given failed: client was not added to subscriber table")
	}
	// Given: we have the connection's client ID in our client ID set
	clientIDSet = &idSet{set: make(map[string]struct{})}
	clientIDSet.set[connection.ClientID] = struct{}{}
	// assert
	if clientIDSet.isClientIDUnique(connection.ClientID) {
		t.Fatalf("given failed: client id should be not be unique because it should have been added")
	}
	// When: there is a disconnect
	HandleDisconnect(connection)
	// Then: the client is not in our subscribers table
	subscriberResult, err := subscribers.getSubscribers(topic)
	if err != nil {
		t.Fatalf("unexpected error %+v\n", err)
	}
	for _, s := range subscriberResult {
		if s.ClientID == connection.ClientID {
			t.Fatalf("expected connection (client) to have been removed from table")
		}
	}
	// Then: the client's client id is not in our client id set
	if !clientIDSet.isClientIDUnique(connection.ClientID) {
		t.Fatalf("client id should be unique")
	}
}
