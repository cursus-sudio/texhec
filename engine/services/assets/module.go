package assets

import (
	"errors"
	"reflect"
)

type AssetModule interface {
	// takes asset struct pointer
	// for each [AssetID] property sets its value to its `path` struct tag value
	InitializeProperties(pointerToStruct any) error
}

type assetModule struct{}

func newAssetModule() AssetModule {
	return &assetModule{}
}

func (a *assetModule) InitializeProperties(pointer any) error {
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
			a.InitializeProperties(fieldValue.Addr().Interface())
			continue
		}

		if fieldType.Type != reflect.TypeFor[AssetID]() {
			continue
		}
		path := fieldType.Tag.Get("path")
		if path == "" {
			continue
		}

		fieldValue.Set(reflect.ValueOf(AssetID(path)))
	}
	return nil
}
