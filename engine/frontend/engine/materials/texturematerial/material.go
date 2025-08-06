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
	Model      int32 `uniform:"model"`
	Camera     int32 `uniform:"camera"`
	Projection int32 `uniform:"projection"`
}

type textureMaterialServices[Projection any] struct {
	window        window.Api
	assetsService assets.Assets
	locations     locations
	getProjection func(Projection) mgl32.Mat4
}

type textureMaterial[Projection any] struct {
	vertSource, fragSource string
	services               *textureMaterialServices[Projection]
	parameters             []program.Parameter
}

func newTextureMaterial[Projection any](
	vertSource, fragSource string,
	window window.Api,
	assetsService assets.Assets,
	parameters []program.Parameter,
	getProjection func(Projection) mgl32.Mat4,
) textureMaterial[Projection] {
	return textureMaterial[Projection]{
		vertSource: vertSource,
		fragSource: fragSource,
		services: &textureMaterialServices[Projection]{
			window:        window,
			assetsService: assetsService,
			getProjection: getProjection,
		},
		parameters: parameters,
	}
}

func (m *textureMaterialServices[Projection]) onFrame(world ecs.World, p program.Program) error {
	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		return err
	}

	var projectionZero Projection
	cameraEntities := world.GetEntitiesWithComponents(ecs.GetComponentType(projectionZero))
	if len(cameraEntities) != 1 {
		return projection.ErrWorldShouldHaveOneProjection
	}

	{
		var projectionComponent Projection
		if err := world.GetComponent(cameraEntities[0], &projectionComponent); err != nil {
			return err
		}
		projection := m.getProjection(projectionComponent)
		gl.UniformMatrix4fv(locations.Projection, 1, false, &projection[0])
	}

	{
		var transformComponent transform.Transform
		if err := world.GetComponent(cameraEntities[0], &transformComponent); err != nil {
			return errors.Join(errors.New("camera misses transform component"), err)
		}

		position := mgl32.Translate3D(
			transformComponent.Pos.X(),
			transformComponent.Pos.Y(),
			transformComponent.Pos.Z(),
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

func (m *textureMaterialServices[Projection]) useForEntity(world ecs.World, p program.Program, entityId ecs.EntityId) error {
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
		transformComponent.Pos.X(),
		transformComponent.Pos.Y(),
		transformComponent.Pos.Z(),
	)

	rotation := transformComponent.Rotation.Mat4()

	scale := mgl32.Scale3D(
		transformComponent.Size.X()/2,
		transformComponent.Size.Y()/2,
		transformComponent.Size.Z()/2,
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

func (m *textureMaterial[Projection]) Material() material.MaterialStorageAsset {
	return material.NewMaterialStorageAsset(
		m.vertSource,
		m.fragSource,
		m.services.onFrame,
		m.services.useForEntity,
		m.parameters,
	)
}
