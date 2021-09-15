package lib

import "testing"

func TestGetFixedHeaderFieldsSuccess(t *testing.T) {
	// Given: A stream of bytes representing a mqtt header
	// When: we try to decode the header
	// Then: we get the values as expected
}

func TestGetRemLengthSuccess(t *testing.T) {
	// Given: some encoded remaining lengths
	// When: we decode them
	// Then: we get the values as expected
}

func TestGetRemLengthSuccessInvalidRemLength(t *testing.T) {
	// Given: some encoded remaining lengths using more than 4 bytes
	// When: we decode them
	// Then: we get an error
}

func TestIsValidUTF8EncodedSuccess(t *testing.T) {
	// Given: some valid utf8 encoded string
	// When: we try to validated
	// Then: we get that it is valid
}


func TestIsValidUTF8EncodedInvalid(t *testing.T) {
	// Given: some invalid utf8 encoded string
	// When: we try to validated
	// Then: we get an error
}