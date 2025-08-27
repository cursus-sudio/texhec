package material

import "frontend/services/assets"

// this component says which materials are used. this allows to cache world effectively
type WorldMeshesAndTextures struct {
	Textures []assets.AssetID
}

func NewWorldMeshesAndTextures(
	textures []assets.AssetID,
) WorldMeshesAndTextures {
	return WorldMeshesAndTextures{
		Textures: textures,
	}
}

// this component says which entities use material
type TextureMaterialComponent struct{}
