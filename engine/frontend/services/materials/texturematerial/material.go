package texturematerial

import (
	"errors"
	"frontend/components/material"
	"frontend/components/projection"
	"frontend/components/texture"
	"frontend/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type locations struct {
	Model      int32 `uniform:"model"`
	Camera     int32 `uniform:"camera"`
	Projection int32 `uniform:"projection"`
}

type textureMaterialServices struct {
	window        window.Api
	assetsService assets.Assets
	locations     locations
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

func (m *textureMaterialServices) onFrame(world ecs.World, p program.Program) error {
	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		return err
	}

	cameraEntities := world.GetEntitiesWithComponents(ecs.GetComponentType(projection.Projection{}))
	if len(cameraEntities) != 1 {
		return projection.ErrWorldShouldHaveOneProjection
	}

	{
		var projectionComponent projection.Projection
		if err := world.GetComponent(cameraEntities[0], &projectionComponent); err != nil {
			return err
		}
		gl.UniformMatrix4fv(locations.Projection, 1, false, &projectionComponent.Projection[0])
	}

	{
		var transformComponent transform.Transform
		if err := world.GetComponent(cameraEntities[0], &transformComponent); err != nil {
			return errors.Join(errors.New("camera misses transform component"), err)
		}

		position := mgl32.Translate3D(
			transformComponent.Pos.X,
			transformComponent.Pos.Y,
			transformComponent.Pos.Z,
		)

		rotation := transformComponent.Rotation.
			Mul(mgl32.QuatRotate(mgl32.DegToRad(180), mgl32.Vec3{0, 1, 0})).
			Mat4()

		matrices := []mgl32.Mat4{
			position,
			rotation,
		}
		var camera mgl32.Mat4
		for i, matrix := range matrices {
			if i == 0 {
				camera = matrix
				continue
			}
			camera = camera.Mul4(matrix)
		}

		gl.UniformMatrix4fv(locations.Camera, 1, false, &camera[0])
	}
	return nil
}

func (m *textureMaterialServices) useForEntity(world ecs.World, p program.Program, entityId ecs.EntityId) error {
	// texture
	var textureComponent texture.Texture
	if err := world.GetComponent(entityId, &textureComponent); err != nil {
		return err
	}
	textureAsset, err := assets.GetAsset[texture.TextureCachedAsset](m.assetsService, textureComponent.ID)
	if err != nil {
		return err
	}

	// transform
	var transformComponent transform.Transform
	if err := world.GetComponent(entityId, &transformComponent); err != nil {
		return err
	}

	// locations
	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		return err
	}

	position := mgl32.Translate3D(
		transformComponent.Pos.X,
		transformComponent.Pos.Y,
		transformComponent.Pos.Z,
	)

	rotation := transformComponent.Rotation.Mat4()

	scale := mgl32.Scale3D(
		transformComponent.Size.X/2,
		transformComponent.Size.Y/2,
		transformComponent.Size.Z/2,
	)

	matrices := []mgl32.Mat4{
		position,
		rotation,
		scale,
	}
	var model mgl32.Mat4
	for i, matrix := range matrices {
		if i == 0 {
			model = matrix
			continue
		}
		model = model.Mul4(matrix)
	}
	gl.UniformMatrix4fv(locations.Model, 1, false, &model[0])

	textureAsset.Texture().Use()
	return nil
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
