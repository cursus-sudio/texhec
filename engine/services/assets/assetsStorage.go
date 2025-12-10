package assets

import (
	"engine/services/httperrors"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type AssetsStorageBuilder interface {
	RegisterAsset(AssetID, func() (any, error))
	RegisterExtension(
		/* shouldn't have dots and be after dots in asset */ extension string,
		loader func(id AssetID) (any, error),
	)
	Build() (AssetsStorage, []error)
}

type assetStorageBuilder struct {
	filePrefix       string
	errs             []error
	assetGetters     map[AssetID]func() (any, error)
	extensionGetters map[string]func(AssetID) (any, error)
}

func NewAssetsStorageBuilder(filePrefix string) AssetsStorageBuilder {
	if len(filePrefix) != 0 && filePrefix[len(filePrefix)-1] != '/' {
		filePrefix += "/"
	}
	return &assetStorageBuilder{
		filePrefix:       filePrefix,
		errs:             []error{},
		assetGetters:     map[AssetID]func() (any, error){},
		extensionGetters: map[string]func(AssetID) (any, error){},
	}
}

func (b *assetStorageBuilder) RegisterAsset(id AssetID, getter func() (any, error)) {
	if _, ok := b.assetGetters[id]; ok {
		err := errors.Join(
			httperrors.Err409,
			fmt.Errorf("\"%s\" id is already registered", id),
		)
		b.errs = append(b.errs, err)
		return
	}
	b.assetGetters[id] = getter
}

func (b *assetStorageBuilder) RegisterExtension(id string, getter func(AssetID) (any, error)) {
	if _, ok := b.extensionGetters[id]; ok {
		err := errors.Join(
			httperrors.Err409,
			fmt.Errorf("\"%s\" extension is already registered", id),
		)
		b.errs = append(b.errs, err)
		return
	}
	b.extensionGetters[id] = getter
}

func (b *assetStorageBuilder) Build() (AssetsStorage, []error) {
	if len(b.errs) != 0 {
		return nil, b.errs
	}
	return &assetsStorage{
		filePrefix:       AssetID(b.filePrefix),
		assetGetters:     b.assetGetters,
		extensionGetters: b.extensionGetters,
	}, nil
}

//

type AssetsStorage interface {
	Get(id AssetID) (any, error)
}

type assetsStorage struct {
	filePrefix       AssetID
	assetGetters     map[AssetID]func() (any, error)
	extensionGetters map[string]func(AssetID) (any, error)
}

func (s *assetsStorage) Get(asset AssetID) (any, error) {
	if getter, ok := s.assetGetters[asset]; ok {
		return getter()
	}
	parts := strings.Split(string(asset), ".")
	extension := parts[len(parts)-1]
	if getter, ok := s.extensionGetters[extension]; ok {
		return getter(s.filePrefix + asset)
	}
	return nil, httperrors.Err404
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
