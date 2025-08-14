package texturematerial

import (
	"errors"
	"frontend/engine/components/material"
	"frontend/engine/components/mesh"
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
	program  func() (program.Program, error)
	services *textureMaterialServices
}

func newTextureMaterial(
	program func() (program.Program, error),
	window window.Api,
	assetsService assets.Assets,
) textureMaterial {
	return textureMaterial{
		program: program,
		services: &textureMaterialServices{
			window:        window,
			assetsService: assetsService,
		},
	}
}

func (m *textureMaterial) Material() material.MaterialStorageAsset {
	return material.NewMaterialStorageAsset(
		m.program,
		m.services.render,
	)
}

//

type textureMaterialServices struct {
	window        window.Api
	assetsService assets.Assets
	locations     locations
}

func (m *textureMaterialServices) render(world ecs.World, p program.Program, entities []ecs.EntityID) error {
	projectionsMvp := map[ecs.ComponentType]mgl32.Mat4{}
	for _, entityId := range entities {
		usedProjection, err := ecs.GetComponent[projection.UsedProjection](world, entityId)
		if err != nil {
			return err
		}

		mvp, ok := projectionsMvp[usedProjection.ProjectionComponent]
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

			cameraTransformComponent, err := ecs.GetComponent[transform.Transform](world, camera)
			if err != nil {
				return errors.Join(errors.New("camera misses transform component"), err)
			}

			projectionMat4 := projectionComponent.Mat4()
			cameraTransformMat4 := projectionComponent.ViewMat4(cameraTransformComponent)

			mvp = projectionMat4.Mul4(cameraTransformMat4)
			projectionsMvp[usedProjection.ProjectionComponent] = mvp
		}

		// texture

		textureComponent, err := ecs.GetComponent[texture.Texture](world, entityId)
		if err != nil {
			continue
		}
		transformComponent, err := ecs.GetComponent[transform.Transform](world, entityId)
		if err != nil {
			continue
		}
		meshComponent, err := ecs.GetComponent[mesh.Mesh](world, entityId)
		if err != nil {
			continue
		}

		textureAsset, err := assets.GetAsset[texture.TextureCachedAsset](m.assetsService, textureComponent.ID)
		if err != nil {
			continue
		}

		// locations
		locations, err := program.GetProgramLocations[locations](p)
		if err != nil {
			continue
		}

		model := transformComponent.Mat4()
		mvp = mvp.Mul4(model)
		gl.UniformMatrix4fv(locations.Mvp, 1, false, &mvp[0])

		textureAsset.Texture().Use()

		// mesh
		meshAsset, err := assets.GetAsset[mesh.MeshCachedAsset](m.assetsService, meshComponent.ID)
		if err != nil {
			return err
		}

		meshAsset.VAO().Draw()
	}
	return nil
}
