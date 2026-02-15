package test

import (
	"engine/modules/assets"
	"testing"
)

type asset struct{ released bool }

func (a *asset) Release() { a.released = true }

//

const assetPath = "asset.png"

func TestAssets(t *testing.T) {
	setup := NewSetup()
	fetched := false
	setup.Extensions.Register("png", func(path assets.Path) (any, error) {
		fetched = true
		return &asset{}, nil
	})
	assetID, ok := setup.Assets.PathID(assetPath)
	if !ok {
		t.Error("registered path extension yet it wan't detected")
		return
	}
	if fetched {
		t.Error("fetched asset prematurely")
		return
	}

	asset, err := assets.GetAsset[*asset](setup.Assets, assetID)
	if err != nil {
		t.Error(err)
		return
	}

	if !fetched {
		t.Error("didn't fetch asset using extension dispatcher")
		return
	}

	if asset.released {
		t.Error("prematurely released asset")
		return
	}

	setup.Assets.Release(assetID)

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
