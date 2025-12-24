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

type CachableAsset interface {
	Cache() (CachedAsset, error)
}

type AssetID string

// add asset struct

type Assets interface {
	Get(AssetID) (any, error)
	Release(...AssetID)
	ReleaseAll()
}

type assets struct {
	assetStorage AssetsStorage
	cachedAssets AssetsCache
}

func NewAssets(
	storage AssetsStorage,
	cached AssetsCache,
) Assets {
	return &assets{
		assetStorage: storage,
		cachedAssets: cached,
	}
}

func (a *assets) Get(id AssetID) (any, error) {
	{
		cached, err := a.cachedAssets.Get(id)
		if err == nil {
			return cached, nil
		}
	}
	stored, err := a.assetStorage.Get(id)
	if err != nil {
		return nil, err
	}
	var cached CachedAsset
	if cachedAsset, ok := stored.(CachedAsset); ok {
		cached = cachedAsset
	} else if cachableAsset, ok := stored.(CachableAsset); ok {
		cached, err = cachableAsset.Cache()
		if err != nil {
			return nil, err
		}
		// } else if cacher, ok := a.cachers[reflect.TypeOf(stored)]; ok {
		// 	cached, err = cacher(stored)
		// 	if err != nil {
		// 		return nil, err
		// 	}
	} else {
		return nil, ErrStoredAssetIsntCachable
	}
	if err := a.cachedAssets.Set(id, cached); err != nil {
		return nil, err
	}
	return cached, nil
}

func (a *assets) Release(ids ...AssetID) {
	for _, id := range ids {
		a.cachedAssets.Delete(id)
	}
}
func (a *assets) ReleaseAll() { a.cachedAssets.DeleteAll() }

var (
	ErrAssetHasDifferentType error = errors.New("asset is not of requested type")
)

func GetAsset[Asset any](assets Assets, assetID AssetID) (Asset, error) {
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
