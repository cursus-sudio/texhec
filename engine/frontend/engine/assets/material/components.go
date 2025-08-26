package material

import "frontend/services/assets"

// this component says which materials are used. this allows to cache world effectively
type WorldTextureMaterialComponent struct {
	Textures []assets.AssetID
	Meshes   []assets.AssetID
}

func NewWorldTextureMaterialComponent(
	textures []assets.AssetID,
	meshes []assets.AssetID,
) WorldTextureMaterialComponent {
	return WorldTextureMaterialComponent{
		Textures: textures,
		Meshes:   meshes,
	}
}

// this component says which entities use material
type TextureMaterialComponent struct{}
