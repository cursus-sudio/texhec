package collider

import (
	"frontend/services/assets"
	"frontend/services/colliders"
)

type Collider struct{ ID assets.AssetID }

func NewCollider(id assets.AssetID) Collider {
	return Collider{ID: id}
}

type ColliderAsset interface {
	assets.GoAsset
	Collider() colliders.Collider
}

type colliderAsset struct {
	assets.GoAsset
	collider colliders.Collider
}

func NewColliderStorageAsset(c colliders.Collider) ColliderAsset {
	asset := &colliderAsset{collider: c}
	asset.GoAsset = assets.NewGoAsset(asset)
	return asset
}

func (a *colliderAsset) Collider() colliders.Collider {
	return a.collider
}
