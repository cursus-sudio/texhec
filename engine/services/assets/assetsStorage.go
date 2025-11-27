package assets

import (
	"engine/services/httperrors"
	"errors"
	"fmt"
	"reflect"
)

type AssetsStorageBuilder interface {
	RegisterAsset(AssetID, func() (any, error))
	Build() (AssetsStorage, []error)
}

type assetStorageBuilder struct {
	errs    []error
	getters map[AssetID]func() (any, error)
}

func NewAssetsStorageBuilder() AssetsStorageBuilder {
	return &assetStorageBuilder{
		errs:    []error{},
		getters: map[AssetID]func() (any, error){},
	}
}

func (b *assetStorageBuilder) RegisterAsset(id AssetID, getter func() (any, error)) {
	if _, ok := b.getters[id]; ok {
		err := errors.Join(
			httperrors.Err409,
			fmt.Errorf("\"%s\" id is already registered", id),
		)
		b.errs = append(b.errs, err)
		return
	}
	b.getters[id] = getter
}

func (b *assetStorageBuilder) Build() (AssetsStorage, []error) {
	if len(b.errs) != 0 {
		return nil, b.errs
	}
	return &assetsStorage{
		getters: b.getters,
	}, nil
}

//

type AssetsStorage interface {
	Get(id AssetID) (any, error)
}

type assetsStorage struct {
	getters map[AssetID]func() (any, error)
}

func (s *assetsStorage) Get(id AssetID) (any, error) {
	getter, ok := s.getters[id]
	if !ok {
		return nil, httperrors.Err404
	}
	return getter()
}

func StorageGet[Asset any](s AssetsStorage, id AssetID) (Asset, error) {
	rawAsset, err := s.Get(id)
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
