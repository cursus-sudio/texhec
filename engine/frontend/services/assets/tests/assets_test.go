package test

import (
	"errors"
	"frontend/services/assets"
	"testing"
)

type storageAsset struct{}

func (a *storageAsset) Cache() (assets.CachedAsset, error) { return &cachedAsset{released: false}, nil }

//

type cachedAsset struct{ released bool }

func (a *cachedAsset) Release() { a.released = true }

//

const notAssetID = "not_asset"
const assetID = "asset"

func TestAssets(t *testing.T) {
	storageBuilder := assets.NewAssetsStorageBuilder()
	fetched := false
	storageBuilder.RegisterAsset(assetID, func() (any, error) {
		fetched = true
		return &storageAsset{}, nil
	})
	storage, errs := storageBuilder.Build()
	if len(errs) != 0 {
		err := errors.Join(errs...)
		t.Error(err)
		return
	}
	if fetched {
		t.Error("fetched asset on build instead of on get")
		return
	}

	cache := assets.NewCachedAssets()
	assets := assets.NewAssets(storage, cache)

	rawAsset, err := assets.Get(assetID)

	if err != nil {
		t.Error(err)
		return
	}

	asset := rawAsset.(*cachedAsset)

	if asset.released {
		t.Error("prematurely released asset")
		return
	}

	assets.Release(assetID)

	if !asset.released {
		t.Error("assets wasn't released")
		return
	}
}

// var cachedAssetInstance = &cachedAsset{released: false}
//
// func TestAssetsStorage(t *testing.T) {
// 	b := assets.NewAssetsStorageBuilder()
// 	fetched := false
// 	b.RegisterAsset(assetID, func() assets.StorageAsset {
// 		fetched = true
// 		return &storageAsset{}
// 	})
// 	storage, errs := b.Build()
// 	if len(errs) != 0 {
// 		err := errors.Join(append([]error{errors.New("error building asset storage")}, errs...)...)
// 		t.Error(err)
// 		return
// 	}
// 	if fetched {
// 		t.Error("fetched asset after building storage and before getting it")
// 		return
// 	}
// 	fetchedAsset, err := storage.Get(assetID)
// 	if err != nil {
// 		t.Error(errors.Join(errors.New("error getting asset which should exist"), err))
// 		return
// 	}
// 	if fetchedAsset != cachedAssetInstance {
// 		t.Errorf("unexpected asset value. expected \"%v\" and got \"%v\"", cachedAssetInstance, fetchedAsset)
// 		return
// 	}
// 	if !fetched {
// 		t.Error("haven't fetched asset after getting it which is interesting because asset value matches")
// 		return
// 	}
// 	_, err = storage.Get(notAssetID)
// 	if !errors.Is(err, httperrors.Err404) {
// 		t.Errorf("expected \"%v\" error. got \"%v\" error", err, httperrors.Err404)
// 		return
// 	}
// }
//
// func TestAssetsCache(t *testing.T) {
// 	cache := assets.NewCachedAssets()
//
// 	err := cache.Set(assetID, cachedAssetInstance)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	cachedAsset, err := cache.Get(assetID)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if cachedAsset != cachedAssetInstance {
// 		t.Errorf("expected cached asset to be \"%v\" but got \"%v\"", cachedAssetInstance, cachedAsset)
// 		return
// 	}
// 	_, err = cache.Get(notAssetID)
// 	if !errors.Is(err, httperrors.Err404) {
// 		t.Errorf("expected \"%s\" error but got \"%s\"", httperrors.Err404, err)
// 		return
// 	}
// 	if cachedAssetInstance.released {
// 		t.Error("cleaned asset before deleting it")
// 		return
// 	}
// 	cache.Delete(assetID)
// 	if !cachedAssetInstance.released {
// 		t.Error("didn't clean asset after deleting it")
// 		return
// 	}
// }
