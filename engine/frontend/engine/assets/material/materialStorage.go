package material

import (
	"frontend/engine/components/material"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/media/window"
	"shared/services/logger"
)

type textureMaterial struct {
	program  func() (program.Program, error)
	services *materialCache
}

func newTextureMaterial(
	program func() (program.Program, error),
	window window.Api,
	assetsStorage assets.AssetsStorage,
	logger logger.Logger,
	entitiesQueryAdditionalArguments []ecs.ComponentType,
) textureMaterial {
	return textureMaterial{
		program: program,
		services: &materialCache{
			window:        window,
			assetsStorage: assetsStorage,
			logger:        logger,

			entitiesQueryAdditionalArguments: entitiesQueryAdditionalArguments,
		},
	}
}

func (m *textureMaterial) Material() material.MaterialStorageAsset {
	return material.NewMaterialStorageAsset(
		m.program,
		m.services.render,
	)
}
