package assets

import (
	"engine/services/httperrors"
	"reflect"
	"sync"
)

type CachedAsset interface {
	Release()
}

type AssetsCache interface {
	Get(id AssetID) (any, error)             // 404
	Set(id AssetID, asset CachedAsset) error // 409
	Delete(id AssetID)
	DeleteAll()
}

type cachedAssets struct {
	mutex   sync.Mutex
	cachers map[reflect.Type]func(any) (CachedAsset, error)
	assets  map[AssetID]CachedAsset
}

func NewCachedAssets() AssetsCache {
	return &cachedAssets{
		mutex:  sync.Mutex{},
		assets: map[AssetID]CachedAsset{},
	}
}

func (assets *cachedAssets) Get(id AssetID) (any, error) {
	assets.mutex.Lock()
	defer assets.mutex.Unlock()
	asset, ok := assets.assets[id]
	if !ok {
		return nil, httperrors.Err404
	}
	return asset, nil
}

func (assets *cachedAssets) Set(id AssetID, asset CachedAsset) error {
	assets.mutex.Lock()
	defer assets.mutex.Unlock()
	if _, ok := assets.assets[id]; ok {
		return httperrors.Err409
	}
	assets.assets[id] = asset
	return nil
}

// if asset implements ReleasableAsset interface Release is called
func (assets *cachedAssets) Delete(id AssetID) {
	assets.mutex.Lock()
	defer assets.mutex.Unlock()
	asset, ok := assets.assets[id]
	if !ok {
		return
	}
	releasableAsset, ok := asset.(CachedAsset)
	if ok {
		releasableAsset.Release()
	}
	delete(assets.assets, id)
}

func (assets *cachedAssets) DeleteAll() {
	for key, asset := range assets.assets {
		delete(assets.assets, key)
		asset.Release()
	}
}
