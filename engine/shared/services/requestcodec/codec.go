package requestcodec

import (
	"encoding/json"
	"errors"
	"reflect"
)

type name string

type Encoded struct {
	Name  name   `json:"name"`
	Bytes []byte `json:"bytes"`
}

type Codec interface {
	// panics when type is not registered
	Encode(any) []byte
	// can return error when its not encodable or type is not registered
	TryEncode(any) ([]byte, error)

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

func (codec *codec) Encode(model any) []byte {
	bytes, err := codec.TryEncode(model)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (codec *codec) TryEncode(model any) ([]byte, error) {
	modelType := reflect.TypeOf(model)
	name, ok := codec.namesByType[modelType]
	if !ok {
		return nil, ErrTypeIsNotRegistered
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
	model := reflect.Zero(modelType).Interface()
	if err := json.Unmarshal(encoded.Bytes, &model); err != nil {
		return nil, errors.Join(ErrInvalidBytes, err)
	}
	return model, nil
}
