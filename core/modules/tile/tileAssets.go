package tile

import (
	"core/modules/definition"
	"frontend/services/assets"
	"shared/services/datastructures"
)

type TileAssets interface {
	AddType(addedAssets datastructures.SparseArray[definition.DefinitionID, assets.AssetID])
}
