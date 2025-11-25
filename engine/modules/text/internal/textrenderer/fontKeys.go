package textrenderer

import (
	"engine/services/assets"
	"sync"
)

type FontKeys interface {
	GetKey(assets.AssetID) FontKey
	GetAsset(FontKey) (assets.AssetID, bool)
}

type fontKeys struct {
	fontsKeys map[assets.AssetID]FontKey
	keysFonts []*assets.AssetID
	mutex     sync.Mutex
	i         FontKey
}

func NewFontKeys() FontKeys {
	return &fontKeys{
		fontsKeys: make(map[assets.AssetID]FontKey),
		keysFonts: make([]*assets.AssetID, 0),
		mutex:     sync.Mutex{},
		i:         FontKey(0),
	}
}

func (k *fontKeys) GetKey(asset assets.AssetID) FontKey {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	fontKey, ok := k.fontsKeys[asset]
	if ok {
		return fontKey
	}

	k.i += 1

	fontKey = k.i
	k.fontsKeys[asset] = fontKey
	for int(fontKey) >= len(k.keysFonts) {
		k.keysFonts = append(k.keysFonts, nil)
	}
	k.keysFonts[fontKey] = &asset

	return fontKey
}

func (k *fontKeys) GetAsset(key FontKey) (assets.AssetID, bool) {
	if int(key) >= len(k.keysFonts) {
		return assets.AssetID(""), false
	}
	asset := k.keysFonts[key]
	if asset == nil {
		return assets.AssetID(""), false
	}

	return *asset, true
}
