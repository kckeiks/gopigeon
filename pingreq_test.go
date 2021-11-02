package gopigeon

import (
	"reflect"
	"testing"
)

func TestHandlePingreqSuccess(t *testing.T) {
	// Given: ping request
	// Note: Fixed hdr reader already Read  all the data before calling pingreq handler function
	fh, _ := newPingreqReqRequest()
	conn := newTestClient([]byte{})
	// When: we send it
	err := HandlePingreq(conn, fh)
	// Then: we get a pingres
	if err != nil {
		t.Fatalf("unexpected error")
	}
	expectedResult := newEncodedPingres()
	result := make([]byte, len(expectedResult))
	conn.Conn.Read(result)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("got encoded PingRes %b but expected %b,", result, expectedResult)
	}
}
