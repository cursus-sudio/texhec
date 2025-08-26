package material

import (
	"frontend/engine/components/material"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/media/window"
)

type textureMaterial struct {
	program  func() (program.Program, error)
	services *materialCache
}

func newTextureMaterial(
	program func() (program.Program, error),
	window window.Api,
	assetsStorage assets.AssetsStorage,
	console console.Console,
	entitiesQueryAdditionalArguments []ecs.ComponentType,
) textureMaterial {
	return textureMaterial{
		program: program,
		services: &materialCache{
			window:        window,
			assetsStorage: assetsStorage,
			console:       console,

			entitiesQueryAdditionalArguments: entitiesQueryAdditionalArguments,

			materialWorldCache: &materialWorldCache{},
		},
	}
}

func (m *textureMaterial) Material() material.MaterialStorageAsset {
	return material.NewMaterialStorageAsset(
		m.program,
		m.services.render,
	)
}
