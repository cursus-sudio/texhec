package internal

import (
	"engine/modules/registry"
	"engine/services/ecs"
	"engine/services/logger"
	"errors"
	"fmt"
	"reflect"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	Logger      logger.Logger `inject:"1"`
	tags        []string
	handlers    []func(string) ecs.EntityID
	presentTags map[string]any
}

func NewService(c ioc.Dic) registry.Service {
	return &service{
		Logger:      ioc.Get[logger.Logger](c),
		tags:        nil,
		handlers:    nil,
		presentTags: make(map[string]any),
	}
}

func (s *service) Register(structTagKey string, handler func(structTagValue string) ecs.EntityID) {
	if _, ok := s.presentTags[structTagKey]; ok {
		s.Logger.Warn(errors.Join(
			fmt.Errorf("already registered struct tag key"),
			fmt.Errorf("struct tag is already registered \"%v\"", structTagKey),
		))
		return
	}

	s.presentTags[structTagKey] = nil
	s.tags = append(s.tags, structTagKey)
	s.handlers = append(s.handlers, handler)
}

func (s *service) populateValue(v reflect.Value) error {
	t := v.Type()
	var err error

	for i := range t.NumField() {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if fieldType.Type != reflect.TypeFor[ecs.EntityID]() {
			if fieldType.Type.Kind() == reflect.Struct {
				if e := s.populateValue(fieldValue); e != nil {
					err = e
				}
			}
			continue
		}
		for tagIndex, tagName := range s.tags {
			tagValue, ok := fieldType.Tag.Lookup(tagName)
			if !ok {
				continue
			}
			tagHandler := s.handlers[tagIndex]
			entity := tagHandler(tagValue)
			fieldValue.Set(reflect.ValueOf(entity))
		}
	}

	return err
}

func (s *service) Populate(structPointer any) error {
	v := reflect.ValueOf(structPointer)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return registry.ErrExpectedPointerToAStruct
	}

	return s.populateValue(v.Elem())
}
