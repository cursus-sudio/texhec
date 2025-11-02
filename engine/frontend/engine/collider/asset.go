package collider

import (
	"frontend/services/assets"

	"github.com/go-gl/mathgl/mgl32"
)

type RangeTarget bool

const (
	Leaf   RangeTarget = false
	Branch RangeTarget = true
)

type Range struct {
	Target       RangeTarget
	First, Count uint32
}

func NewRange(target RangeTarget, first, count uint32) Range {
	return Range{Target: target, First: first, Count: count}
}

// todo add normals and store aabb
type Polygon struct{ A, B, C mgl32.Vec3 }

func NewPolygon(a, b, c mgl32.Vec3) Polygon {
	return Polygon{a, b, c}
}

type ColliderAsset interface {
	assets.GoAsset
	// first aabb is the entry point
	AABBs() []AABB
	// []Range element index corresponds to []AABB element index
	Ranges() []Range

	Polygons() []Polygon
}

type colliderAsset struct {
	// assets.GoAsset
	aabbs    []AABB
	ranges   []Range
	polygons []Polygon
}

func NewColliderStorageAsset(
	aabbs []AABB,
	ranges []Range,
	polygons []Polygon,
) ColliderAsset {
	asset := &colliderAsset{
		aabbs:    aabbs,
		ranges:   ranges,
		polygons: polygons,
	}
	// asset.GoAsset = assets.NewGoAsset(asset)
	return asset
}

func (a *colliderAsset) AABBs() []AABB       { return a.aabbs }
func (a *colliderAsset) Ranges() []Range     { return a.ranges }
func (a *colliderAsset) Polygons() []Polygon { return a.polygons }
func (a *colliderAsset) Release()            {}
