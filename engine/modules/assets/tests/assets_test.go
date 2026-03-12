package test

import (
	"engine/modules/assets"
	"engine/modules/registry"
	"engine/services/ecs"
	"testing"
)

type asset struct{ released bool }

func (a *asset) Release() { a.released = true }

//

type Definitions struct {
	Asset ecs.EntityID `path:"asset.png"`
}

func TestAssets(t *testing.T) {
	setup := NewSetup()
	fetched := false
	setup.Assets.Register("png", func(path assets.PathComponent) (assets.Asset, error) {
		fetched = true
		return &asset{}, nil
	})
	definitions, err := registry.GetRegistry[Definitions](setup.Registry)
	if err != nil {
		t.Error("registered path extension yet it wan't detected")
		return
	}
	if fetched {
		t.Error("fetched asset prematurely")
		return
	}

	asset, err := assets.GetAsset[*asset](setup.Assets, definitions.Asset)
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

	setup.Assets.Release(definitions.Asset)

	if !asset.released {
		t.Error("assets wasn't released")
		return
	}
}
