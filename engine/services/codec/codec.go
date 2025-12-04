package codec

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
)

type Codec interface {
	Encode(any) ([]byte, error)

	// can return:
	// ErrInvalidInput
	Decode([]byte) (any, error)
}

type codec struct{}

func NewCodec(types []reflect.Type) Codec {
	for _, codecType := range types {
		gob.Register(reflect.New(codecType).Interface())
	}
	return &codec{}
}

func (codec *codec) Encode(model any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(&model); err != nil {
		return nil, errors.Join(ErrCannotEncodeType, err)
	}
	return buffer.Bytes(), nil
}

func (codec *codec) Decode(bytesToDecode []byte) (any, error) {
	var value any
	if err := gob.
		NewDecoder(bytes.NewBuffer(bytesToDecode)).
		Decode(&value); err != nil {
		return nil, errors.Join(ErrInvalidBytes, err)
	}
	if value == nil {
		return nil, errors.Join(ErrTypeIsNotRegistered)
	}
	return value, nil
}
