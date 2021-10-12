package gopigeon

import (
	"bytes"
	"reflect"
	"testing"
)

func TestGetFixedHeaderFieldsSuccess(t *testing.T) {
	// Given: A stream of bytes representing a mqtt header
	expectedResult := FixedHeader{PktType: Connect, RemLength: 12}
	h := newTestEncodedFixedHeader(expectedResult)
	// When: we try to decode the header
	result, _ := ReadFixedHeader(bytes.NewBuffer(h))
	// Then: we get the values as expected
	if !reflect.DeepEqual(result, &expectedResult) {
		t.Fatalf("Got FixedHeader %+v but expected %+v,", result, &expectedResult)
	}
}

func TestGetRemLengthEdgeCasesSuccess(t *testing.T) {
	// Given: some encoded remaining lengths
	oneByteMin := []byte{0x00}
	oneByteMax := []byte{0x7F}
	twoByteMin := []byte{0x80, 0x01}
	twoByteMax := []byte{0xFF, 0x7F}
	threeByteMin := []byte{0x80, 0x80, 0x01}
	threeByteMax := []byte{0xFF, 0xFF, 0x7F}
	fourByteMin := []byte{0x80, 0x80, 0x80, 0x01}
	fourByteMax := []byte{0xFF, 0xFF, 0xFF, 0x7F}
	// When: we decode them
	result, _ := ReadRemLength(bytes.NewBuffer(oneByteMin))
	// Then: we get the values as expected
	expectedResult := uint32(0)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
	// When: we decode them
	result, _ = ReadRemLength(bytes.NewBuffer(oneByteMax))
	// Then: we get the values as expected
	expectedResult = uint32(127)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
	// When: we decode them
	result, _ = ReadRemLength(bytes.NewBuffer(twoByteMin))
	// Then: we get the values as expected
	expectedResult = uint32(128)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
	// When: we decode them
	result, _ = ReadRemLength(bytes.NewBuffer(twoByteMax))
	// Then: we get the values as expected
	expectedResult = uint32(16383)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
	// When: we decode them
	result, _ = ReadRemLength(bytes.NewBuffer(threeByteMin))
	// Then: we get the values as expected
	expectedResult = uint32(16384)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
	// When: we decode them
	result, _ = ReadRemLength(bytes.NewBuffer(threeByteMax))
	// Then: we get the values as expected
	expectedResult = uint32(2097151)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
	// When: we decode them
	result, _ = ReadRemLength(bytes.NewBuffer(fourByteMin))
	// Then: we get the values as expected
	expectedResult = uint32(2097152)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
	// When: we decode them
	result, _ = ReadRemLength(bytes.NewBuffer(fourByteMax))
	// Then: we get the values as expected
	expectedResult = uint32(268435455)
	if result != expectedResult {
		t.Fatalf("Got rem length %d but expected %d,", result, expectedResult)
	}
}

func TestGetRemLengthSuccessInvalidRemLength(t *testing.T) {
	// Given: some encoded remaining lengths using more than 4 bytes
	invalidEncodedRemLen := []byte{0x80, 0x80, 0x80, 0x80, 0x01}
	// When: we decode them
	_, err := ReadRemLength(bytes.NewBuffer(invalidEncodedRemLen))
	// Then: we get an error
	if err == nil {
		t.Fatalf("Did not get an error.")
	}
}

func TestIsValidUTF8EncodedSuccess(t *testing.T) {
	// Given: some valid utf8 encoded string
	validEncodedStr := []byte{77, 81, 84, 84} // MQTT
	// When: we try to validated
	result := IsValidUTF8Encoded(validEncodedStr)
	// Then: we get that it is valid
	if !result {
		t.Fatalf("Got IsValidUTF8Encoded false but expected true.")
	}
}

func TestIsValidUTF8EncodedInvalid(t *testing.T) {
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result := IsValidUTF8Encoded([]byte{0xED, 0xA0, 0x80}) // U+D800
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result = IsValidUTF8Encoded([]byte{0xED, 0xBF, 0xBF}) // U+DFFF
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result = IsValidUTF8Encoded([]byte{0x00}) // null character
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result = IsValidUTF8Encoded([]byte{0x01}) // ctrl character
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result = IsValidUTF8Encoded([]byte{0x1F}) // ctrl character
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result = IsValidUTF8Encoded([]byte{0x7F}) // ctrl character
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result = IsValidUTF8Encoded([]byte{0xC2, 0x9F}) // ctrl character
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	result = IsValidUTF8Encoded([]byte{0xEF, 0xBF, 0xBF}) // non-character
	// Then: we get an error
	if result {
		t.Fatalf("Got IsValidUTF8Encoded true but expected false.")
	}
}

func TestReadEncodedStr(t *testing.T) {
	// Given: an encoded string prepended with it's length
	encodedStr := bytes.NewBuffer([]byte{0, 4, 77, 81, 84, 84}) // MQTT
	// When: we decoded it
	result, err := ReadEncodedStr(encodedStr)
	// Then: we get a str representation of it
	if err != nil {
		t.Fatalf("Got an unexpected error.")
	}
	if result != "MQTT" {
		t.Fatalf("Got %s but expected MQTT.", result)
	}
}
