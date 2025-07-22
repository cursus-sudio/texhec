package assets

import (
	"errors"
	"fmt"
	"shared/utils/httperrors"
)

type StorageAsset interface {
	Cache() (CachedAsset, error)
}

type AssetsStorageBuilder interface {
	RegisterAsset(AssetID, func() (StorageAsset, error))
	Build() (AssetsStorage, []error)
}

type assetStorageBuilder struct {
	errs    []error
	getters map[AssetID]func() (StorageAsset, error)
}

func NewAssetsStorageBuilder() AssetsStorageBuilder {
	return &assetStorageBuilder{
		errs:    []error{},
		getters: map[AssetID]func() (StorageAsset, error){},
	}
}

func (b *assetStorageBuilder) RegisterAsset(id AssetID, getter func() (StorageAsset, error)) {
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
	Get(id AssetID) (StorageAsset, error)
}

type assetsStorage struct {
	getters map[AssetID]func() (StorageAsset, error)
}

func (s *assetsStorage) Get(id AssetID) (StorageAsset, error) {
	getter, ok := s.getters[id]
	if !ok {
		return nil, httperrors.Err404
	}
	return getter()
}
