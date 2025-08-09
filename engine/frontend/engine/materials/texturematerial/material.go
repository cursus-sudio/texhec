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

type textureMaterial[Projection projection.Projection] struct {
	vertSource, fragSource string
	services               *textureMaterialServices[Projection]
	parameters             []program.Parameter
}

func newTextureMaterial[Projection projection.Projection](
	vertSource, fragSource string,
	window window.Api,
	assetsService assets.Assets,
	parameters []program.Parameter,
) textureMaterial[Projection] {
	return textureMaterial[Projection]{
		vertSource: vertSource,
		fragSource: fragSource,
		services: &textureMaterialServices[Projection]{
			window:        window,
			assetsService: assetsService,
		},
		parameters: parameters,
	}
}

func (m *textureMaterial[Projection]) Material() material.MaterialStorageAsset {
	return material.NewMaterialStorageAsset(
		m.vertSource,
		m.fragSource,
		m.services.onFrame,
		m.services.useForEntity,
		m.parameters,
	)
}

//

type textureMaterialServices[Projection projection.Projection] struct {
	window        window.Api
	assetsService assets.Assets
	locations     locations

	mvp *mgl32.Mat4
}

func (m *textureMaterialServices[Projection]) onFrame(world ecs.World, _ program.Program) error {
	var projectionZero Projection
	cameraEntities := world.GetEntitiesWithComponents(ecs.GetComponentType(projectionZero))
	if len(cameraEntities) != 1 {
		return projection.ErrWorldShouldHaveOneProjection
	}
	camera := cameraEntities[0]

	var projectionComponent Projection
	if err := world.GetComponents(camera, &projectionComponent); err != nil {
		return err
	}

	var cameraTransformComponent transform.Transform
	if err := world.GetComponents(cameraEntities[0], &cameraTransformComponent); err != nil {
		return errors.Join(errors.New("camera misses transform component"), err)
	}

	projectionMat4 := projectionComponent.Mat4()
	cameraTransformMat4 := projectionComponent.ViewMat4(cameraTransformComponent)

	mvp := projectionMat4.Mul4(cameraTransformMat4)
	m.mvp = &mvp

	return nil
}

func (m *textureMaterialServices[Projection]) useForEntity(world ecs.World, p program.Program, entityId ecs.EntityId) error {
	if m.mvp == nil {
		return material.ErrHaveToCallOnFrame
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
	mvp := m.mvp.Mul4(model)
	gl.UniformMatrix4fv(locations.Mvp, 1, false, &mvp[0])

	textureAsset.Texture().Use()
	return nil
}
