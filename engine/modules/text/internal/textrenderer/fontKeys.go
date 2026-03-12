package textrenderer

import (
	"engine/services/ecs"
	"sync"
)

type FontKeys interface {
	GetKey(ecs.EntityID) FontKey
	GetAsset(FontKey) (ecs.EntityID, bool)
}

type fontKeys struct {
	fontsKeys map[ecs.EntityID]FontKey
	keysFonts []*ecs.EntityID
	mutex     sync.Mutex
	i         FontKey
}

func NewFontKeys() FontKeys {
	return &fontKeys{
		fontsKeys: make(map[ecs.EntityID]FontKey),
		keysFonts: make([]*ecs.EntityID, 0),
		mutex:     sync.Mutex{},
		i:         FontKey(0),
	}
}

func (k *fontKeys) GetKey(asset ecs.EntityID) FontKey {
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

func (k *fontKeys) GetAsset(key FontKey) (ecs.EntityID, bool) {
	if int(key) >= len(k.keysFonts) {
		return 0, false
	}
	asset := k.keysFonts[key]
	if asset == nil {
		return 0, false
	}

	return *asset, true
}
