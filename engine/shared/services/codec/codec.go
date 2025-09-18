package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type name string

type Encoded struct {
	Name  name            `json:"t"`
	Bytes json.RawMessage `json:"v"`
}

type Codec interface {
	Encode(any) ([]byte, error)

	// can return:
	// ErrInvalidInput
	Decode([]byte) (any, error)
}

type codec struct {
	typesByName map[name]reflect.Type
	namesByType map[reflect.Type]name
}

func NewCodec(types []reflect.Type) Codec {
	typesByName := make(map[name]reflect.Type, len(types))
	namesByType := make(map[reflect.Type]name, len(types))
	for _, codecType := range types {
		typeName := name(codecType.String())

		typesByName[typeName] = codecType
		namesByType[codecType] = typeName
	}
	return &codec{
		typesByName: typesByName,
		namesByType: namesByType,
	}
}

func (codec *codec) Encode(model any) ([]byte, error) {
	modelType := reflect.TypeOf(model)
	name, ok := codec.namesByType[modelType]
	if !ok {
		return nil, errors.Join(
			ErrTypeIsNotRegistered,
			fmt.Errorf("codec is missing \"%v\" type", modelType.String()),
		)
	}
	modelBytes, err := json.Marshal(model)
	if err != nil {
		return nil, errors.Join(ErrCannotEncodeType, err)
	}
	encoded := Encoded{
		Name:  name,
		Bytes: modelBytes,
	}
	bytes, err := json.Marshal(encoded)
	if err != nil {
		return nil, errors.Join(ErrCannotEncodeType, err)
	}
	return bytes, nil
}

func (codec *codec) Decode(bytes []byte) (any, error) {
	var encoded Encoded
	if err := json.Unmarshal(bytes, &encoded); err != nil {
		return nil, errors.Join(ErrInvalidBytes, err)
	}
	modelType, ok := codec.typesByName[encoded.Name]
	if !ok {
		return nil, ErrTypeIsNotRegistered
	}
	modelPtr := reflect.New(modelType)
	if err := json.Unmarshal(encoded.Bytes, modelPtr.Interface()); err != nil {
		return nil, errors.Join(ErrInvalidBytes, err)
	}
	model := modelPtr.Elem().Interface()
	return model, nil
}
