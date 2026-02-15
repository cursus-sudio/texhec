package textrenderer

import (
	"engine/modules/assets"
	"sync"
)

type FontKeys interface {
	GetKey(assets.ID) FontKey
	GetAsset(FontKey) (assets.ID, bool)
}

type fontKeys struct {
	fontsKeys map[assets.ID]FontKey
	keysFonts []*assets.ID
	mutex     sync.Mutex
	i         FontKey
}

func NewFontKeys() FontKeys {
	return &fontKeys{
		fontsKeys: make(map[assets.ID]FontKey),
		keysFonts: make([]*assets.ID, 0),
		mutex:     sync.Mutex{},
		i:         FontKey(0),
	}
}

func (k *fontKeys) GetKey(asset assets.ID) FontKey {
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

func (k *fontKeys) GetAsset(key FontKey) (assets.ID, bool) {
	if int(key) >= len(k.keysFonts) {
		return 0, false
	}
	asset := k.keysFonts[key]
	if asset == nil {
		return 0, false
	}

	return *asset, true
}
