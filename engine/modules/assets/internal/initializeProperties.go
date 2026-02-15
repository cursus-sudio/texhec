package internal

import (
	"engine/modules/assets"
	"errors"
	"fmt"
	"reflect"
)

func (a *service) InitializeProperties(pointer any) error {
	pointerValue := reflect.ValueOf(pointer)
	pointerType := pointerValue.Type()
	if pointerType.Kind() != reflect.Pointer {
		return errors.New("instance isn't a pointer")
	}
	structValue := pointerValue.Elem()
	structType := structValue.Type()
	if structType.Kind() != reflect.Struct {
		return errors.New("instance isn't a pointer to a struct")
	}
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)

		if fieldType.Type.Kind() == reflect.Struct {
			if err := a.InitializeProperties(fieldValue.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if fieldType.Type != reflect.TypeFor[assets.ID]() {
			continue
		}
		pathString := fieldType.Tag.Get("path")
		if pathString == "" {
			continue
		}

		path := assets.Path(pathString)
		id, ok := a.PathID(path)
		if !ok {
			extension := a.Extensions.PathExntesion(path)
			return fmt.Errorf("extension \"%v\" isn't registered", extension)
		}

		fieldValue.Set(reflect.ValueOf(id))
	}
	return nil
}
