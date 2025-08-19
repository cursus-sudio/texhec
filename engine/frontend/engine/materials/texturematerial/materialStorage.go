package texturematerial

import (
	"frontend/engine/components/material"
	"frontend/services/assets"
	"frontend/services/graphics/program"
	"frontend/services/media/window"
)

type textureMaterial struct {
	program  func() (program.Program, error)
	services *textureMaterialServices
}

func newTextureMaterial(
	program func() (program.Program, error),
	window window.Api,
	assetsStorage assets.AssetsStorage,
) textureMaterial {
	return textureMaterial{
		program: program,
		services: &textureMaterialServices{
			window:        window,
			assetsStorage: assetsStorage,
			renderCache:   &renderCache{},
		},
	}
}

func (m *textureMaterial) Material() material.MaterialStorageAsset {
	return material.NewMaterialStorageAsset(
		m.program,
		m.services.render,
	)
}
