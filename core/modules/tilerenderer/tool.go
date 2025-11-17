package tilerenderer

import (
	"core/modules/definition"
	"frontend/services/assets"
	"shared/services/datastructures"
)

type TileTool interface {
	AddType(addedAssets datastructures.SparseArray[definition.DefinitionID, assets.AssetID])
}
