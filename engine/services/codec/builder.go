package codec

import "reflect"

type Builder interface {
	Register(any) Builder
	TryRegister(any) error
	Build() Codec
}

type builder struct {
	types map[reflect.Type]struct{}
}

func NewBuilder() Builder {
	return &builder{types: make(map[reflect.Type]struct{})}
}

func (b *builder) Register(codecExample any) Builder {
	if err := b.TryRegister(codecExample); err != nil {
		panic(err)
	}
	return b
}

func (b *builder) TryRegister(codecExample any) error {
	codecType := reflect.TypeOf(codecExample)
	_, ok := b.types[codecType]
	if ok {
		return ErrTypeIsAlreadyRegistered
	}
	b.types[codecType] = struct{}{}
	return nil
}

func (b *builder) Build() Codec {
	types := make([]reflect.Type, 0, len(b.types))
	for codecType := range b.types {
		types = append(types, codecType)
	}
	return NewCodec(types)
}
