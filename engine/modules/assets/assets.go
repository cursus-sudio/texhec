package assets

import (
	"engine/services/httperrors"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrStoredAssetIsntCachable error = errors.New("stored asset isn't cachable")
)

type ID uint32
type Path string
type Asset interface{ Release() }

// add asset struct

type Extensions interface {
	Register(
		/* shouldn't have dots and be after dots in asset */ extension string,
		dispatcher func(path Path) (any, error),
	)
	PathExntesion(Path) string
	ExtensionDispatcher(extension string) (func(Path) (any, error), bool)
}

type Service interface {
	// takes asset struct pointer
	// for each [AssetID] property sets its value to its `path` struct tag value
	InitializeProperties(pointerToStruct any) error
	PathID(Path) (ID, bool)

	// get also caches asset
	Get(ID) (any, error)
	Release(...ID)
	ReleaseAll()
}

var (
	ErrAssetHasDifferentType error = errors.New("asset is not of requested type")
)

func GetAsset[Asset any](assets Service, assetID ID) (Asset, error) {
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
