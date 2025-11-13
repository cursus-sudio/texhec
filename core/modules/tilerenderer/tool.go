package tilerenderer

import (
	"frontend/services/assets"
	"shared/services/datastructures"
)

type TileTool interface {
	AddType(addedAssets datastructures.SparseArray[uint32, assets.AssetID])
}
