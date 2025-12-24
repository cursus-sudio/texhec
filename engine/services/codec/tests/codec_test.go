package test

import (
	"reflect"
	"testing"
)

func TestCodec(t *testing.T) {
	setup := NewSetup()
	value := Type{Value: 48}
	valueBytes, err := setup.codec.Encode(value)
	if err != nil {
		t.Error(err)
		return
	}
	decoded, err := setup.codec.Decode(valueBytes)
	if err != nil {
		t.Error(err)
		return
	}
	valueT := reflect.TypeOf(value)
	decodedT := reflect.TypeOf(decoded)
	if decodedT != valueT {
		t.Errorf("unexpected types mismatch (original) %v != (decoded) %v", valueT, decodedT)
	}
	if decoded != value {
		t.Errorf("unexpected values mismatch (original) %v != (decoded) %v", value, decoded)
	}
}
