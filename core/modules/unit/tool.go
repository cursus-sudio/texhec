package unit

import (
	"frontend/services/assets"
	"shared/services/datastructures"
)

type UnitTool interface {
	AddType(addedAssets datastructures.SparseArray[uint32, assets.AssetID])
}
