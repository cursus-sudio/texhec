package texturematerial

import (
	"errors"
	"frontend/engine/components/material"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type locations struct {
	Mvp int32 `uniform:"mvp"`
}

type textureMaterial struct {
	vertSource, fragSource string
	services               *textureMaterialServices
	parameters             []program.Parameter
}

func newTextureMaterial(
	vertSource, fragSource string,
	window window.Api,
	assetsService assets.Assets,
	parameters []program.Parameter,
) textureMaterial {
	return textureMaterial{
		vertSource: vertSource,
		fragSource: fragSource,
		services: &textureMaterialServices{
			window:        window,
			assetsService: assetsService,
		},
		parameters: parameters,
	}
}

func (m *textureMaterial) Material() material.MaterialStorageAsset {
	return material.NewMaterialStorageAsset(
		m.vertSource,
		m.fragSource,
		m.services.onFrame,
		m.services.useForEntity,
		m.parameters,
	)
}

//

type textureMaterialServices struct {
	window        window.Api
	assetsService assets.Assets
	locations     locations

	projectionsMvp map[ecs.ComponentType]mgl32.Mat4
}

func (m *textureMaterialServices) onFrame(world ecs.World, _ program.Program) error {
	m.projectionsMvp = map[ecs.ComponentType]mgl32.Mat4{}
	return nil
}

func (m *textureMaterialServices) useForEntity(world ecs.World, p program.Program, entityId ecs.EntityId) error {
	var usedProjection projection.UsedProjection
	if err := world.GetComponents(entityId, &usedProjection); err != nil {
		return err
	}

	mvp, ok := m.projectionsMvp[usedProjection.ProjectionComponent]
	if !ok {
		cameraEntities := world.GetEntitiesWithComponents(usedProjection.ProjectionComponent)
		if len(cameraEntities) != 1 {
			return projection.ErrWorldShouldHaveOneProjection
		}
		camera := cameraEntities[0]

		projectionComponent, err := usedProjection.GetCameraProjection(world, camera)
		if err != nil {
			return err
		}

		var cameraTransformComponent transform.Transform
		if err := world.GetComponents(camera, &cameraTransformComponent); err != nil {
			return errors.Join(errors.New("camera misses transform component"), err)
		}

		projectionMat4 := projectionComponent.Mat4()
		cameraTransformMat4 := projectionComponent.ViewMat4(cameraTransformComponent)

		mvp = projectionMat4.Mul4(cameraTransformMat4)
		m.projectionsMvp[usedProjection.ProjectionComponent] = mvp
	}

	// texture
	var textureComponent texture.Texture
	if err := world.GetComponents(entityId, &textureComponent); err != nil {
		return err
	}
	textureAsset, err := assets.GetAsset[texture.TextureCachedAsset](m.assetsService, textureComponent.ID)
	if err != nil {
		return err
	}

	// transform
	var transformComponent transform.Transform
	if err := world.GetComponents(entityId, &transformComponent); err != nil {
		return err
	}

	// locations
	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		return err
	}

	model := transformComponent.Mat4()
	mvp = mvp.Mul4(model)
	gl.UniformMatrix4fv(locations.Mvp, 1, false, &mvp[0])

	textureAsset.Texture().Use()
	return nil
}
