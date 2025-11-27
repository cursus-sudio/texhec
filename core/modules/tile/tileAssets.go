package tile

import (
	"core/modules/definition"
	"engine/services/assets"
	"engine/services/datastructures"
)

type TileAssets interface {
	AddType(addedAssets datastructures.SparseArray[definition.DefinitionID, assets.AssetID])
}
