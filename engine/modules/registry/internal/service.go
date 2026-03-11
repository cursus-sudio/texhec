package internal

import (
	"engine/modules/registry"
	"engine/services/ecs"
	"reflect"
)

type service struct {
	tags        []string
	handlers    []func(string) ecs.EntityID
	presentTags map[string]any
}

func NewService() registry.Service {
	return &service{
		tags:        nil,
		handlers:    nil,
		presentTags: make(map[string]any),
	}
}

func (s *service) Register(structTagKey string, handler func(structTagValue string) ecs.EntityID) error {
	if _, ok := s.presentTags[structTagKey]; ok {
		return registry.ErrAlreadyRegistered
	}

	s.presentTags[structTagKey] = nil
	s.tags = append(s.tags, structTagKey)
	s.handlers = append(s.handlers, handler)
	return nil
}

func (s *service) Populate(structPointer any) error {
	v := reflect.ValueOf(structPointer)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return registry.ErrExpectedPointerToAStruct
	}

	v = v.Elem()
	t := v.Type()

	var err error

fieldLoop:
	for i := range t.NumField() {
		fieldType := t.Field(i)
		if fieldType.Type != reflect.TypeFor[ecs.EntityID]() {
			continue
		}
		fieldValue := v.Field(i)
		for tagIndex, tagName := range s.tags {
			tagValue, ok := fieldType.Tag.Lookup(tagName)
			if !ok {
				continue
			}
			tagHandler := s.handlers[tagIndex]
			entity := tagHandler(tagValue)
			fieldValue.Set(reflect.ValueOf(entity))
			continue fieldLoop
		}
		err = registry.ErrNotFoundHandlerForAField
	}

	return err
}
