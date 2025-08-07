package assets

import (
	"errors"
	"fmt"
	"reflect"
	"shared/utils/httperrors"
)

type AssetID string

type Assets interface {
	Get(AssetID) (any, error)
	Release(...AssetID)
}

type assets struct {
	assetStorage AssetsStorage
	cachedAssets CachedAssets
}

func NewAssets(
	storage AssetsStorage,
	cached CachedAssets,
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
	cached, err := stored.Cache()
	if err != nil {
		return nil, err
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

var (
	ErrAssetHasDifferentType error = errors.New("asset if not of requested type")
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
