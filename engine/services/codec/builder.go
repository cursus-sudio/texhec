package codec

import (
	"engine/services/logger"
	"reflect"
)

type Builder interface {
	Register(any) Builder
	Build() Codec
}

type builder struct {
	logger logger.Logger
	types  map[reflect.Type]struct{}
}

func NewBuilder(logger logger.Logger) Builder {
	return &builder{
		logger: logger,
		types:  make(map[reflect.Type]struct{}),
	}
}

type GobTypesHook interface { // types to register
	GobTypes() []any
}

func (b *builder) Register(codecExample any) Builder {
	codecType := reflect.TypeOf(codecExample)
	if _, ok := b.types[codecType]; ok {
		return b
	}
	b.types[codecType] = struct{}{}

	if h, ok := codecExample.(GobTypesHook); ok {
		for _, t := range h.GobTypes() {
			b.Register(t)
		}
	}
	return b
}

func (b *builder) Build() Codec {
	types := make([]reflect.Type, 0, len(b.types))
	for codecType := range b.types {
		types = append(types, codecType)
	}
	return newCodec(b.logger, types)
}
