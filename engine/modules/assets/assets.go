package assets

import (
	"engine/services/ecs"
	"engine/services/httperrors"
	"errors"
	"fmt"
	"reflect"
)

type Asset interface{ Release() }

// components

type PathComponent struct{ Path string }
type CacheComponent struct{ Cache Asset }

func NewPath(path string) PathComponent   { return PathComponent{path} }
func NewCache(cache Asset) CacheComponent { return CacheComponent{cache} }

// add asset struct

type Service interface {
	Path() ecs.ComponentsArray[PathComponent]
	Cache() ecs.ComponentsArray[CacheComponent]

	Register(
		/* shouldn't have dots and be after dots in asset */ extension string,
		dispatcher func(path PathComponent) (Asset, error),
	)

	// get also caches asset
	Get(ecs.EntityID) (Asset, error)
	Release(...ecs.EntityID)
	ReleaseAll()
}

var (
	ErrAssetHasDifferentType error = errors.New("asset is not of requested type")
	ErrAssetNotFound         error = errors.New("asset not found")
)

func GetAsset[Asset any](assets Service, assetID ecs.EntityID) (Asset, error) {
	rawAsset, err := assets.Get(assetID)
	if err != nil {
		var a Asset
		return a, err
	}
	asset, ok := rawAsset.(Asset)
	if !ok {
		var a Asset
		err := errors.Join(
			httperrors.Err400,
			ErrAssetHasDifferentType,
			fmt.Errorf(
				"asset is of type \"%s\" and expected to be \"%s\"",
				reflect.TypeOf(rawAsset).String(),
				reflect.TypeFor[Asset]().String(),
			),
		)
		return a, err
	}
	return asset, nil
}
