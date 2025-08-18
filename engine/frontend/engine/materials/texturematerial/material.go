package texturematerial

import (
	"errors"
	"frontend/engine/components/material"
	meshcomponent "frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/graphics/texture"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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

//

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
	assetsStorage assets.AssetsStorage,
) textureMaterial {
	return textureMaterial{
		program: program,
		services: &textureMaterialServices{
			window:        window,
			assetsStorage: assetsStorage,
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
	assetsStorage assets.AssetsStorage
	locations     locations

	cachedForWorld    ecs.World
	entitiesLiveQuery ecs.LiveQuery
	meshes            map[assets.AssetID]vao.VAO
	textures          map[assets.AssetID]texture.Texture
}

func (m *textureMaterialServices) cleanUp() {
	for _, vao := range m.meshes {
		vao.Release()
	}
	m.meshes = map[assets.AssetID]vao.VAO{}

	for _, texture := range m.textures {
		texture.Release()
	}
	m.textures = map[assets.AssetID]texture.Texture{}
}

func (m *textureMaterialServices) init(world ecs.World) {
	m.cachedForWorld = world
	m.entitiesLiveQuery = world.QueryEntitiesWithComponents(
		ecs.GetComponentType(TextureMaterialComponent{}),
		ecs.GetComponentType(meshcomponent.Mesh{}),
		ecs.GetComponentType(projection.UsedProjection{}),
		ecs.GetComponentType(transform.Transform{}),
		ecs.GetComponentType(texturecomponent.Texture{}),
	)
	m.meshes = map[assets.AssetID]vao.VAO{}
	m.textures = map[assets.AssetID]texture.Texture{}

	var worldComponent WorldTextureMaterialComponent
	query := world.QueryEntitiesWithComponents(ecs.GetComponentType(WorldTextureMaterialComponent{}))
	for _, entity := range query.Entities() {
		c, err := ecs.GetComponent[WorldTextureMaterialComponent](world, entity)
		if err != nil {
			continue
		}
		worldComponent.Meshes = append(worldComponent.Meshes, c.Meshes...)
		worldComponent.Textures = append(worldComponent.Textures, c.Textures...)
	}

	for _, assetID := range worldComponent.Textures {
		textureAsset, err := assets.StorageGet[texturecomponent.TextureStorageAsset](m.assetsStorage, assetID)
		if err != nil {
			continue
		}

		t, err := texture.NewTexture(textureAsset.Image())
		if err != nil {
			continue
		}

		m.textures[assetID] = t
	}

	for _, assetID := range worldComponent.Meshes {
		meshAsset, err := assets.StorageGet[meshcomponent.MeshStorageAsset](m.assetsStorage, assetID)
		if err != nil {
			continue
		}

		vbo := vbo.NewVBO()
		vbo.SetVertices(meshAsset.Verticies())

		ebo := ebo.NewEBO()
		ebo.SetIndices(meshAsset.Indicies())

		vao := vao.NewVAO(vbo, ebo)
		m.meshes[assetID] = vao
	}
}

func (m *textureMaterialServices) render(world ecs.World, p program.Program) error {
	reCache := m.cachedForWorld != world
	if reCache {
		m.cleanUp()
		m.init(world)
	}
	entities := m.entitiesLiveQuery.Entities()
	if true {
		projectionsMvp := map[ecs.ComponentType]mgl32.Mat4{}

		for _, entity := range entities {
			usedProjection, err := ecs.GetComponent[projection.UsedProjection](world, entity)
			if err != nil {
				return err
			}

			mvp, ok := projectionsMvp[usedProjection.ProjectionComponent]
			if !ok {
				query := world.QueryEntitiesWithComponents(usedProjection.ProjectionComponent)
				cameraEntities := query.Entities()
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

			transformComponent, err := ecs.GetComponent[transform.Transform](world, entity)
			if err != nil {
				continue
			}

			textureComponent, err := ecs.GetComponent[texturecomponent.Texture](world, entity)
			if err != nil {
				continue
			}
			textureAsset, ok := m.textures[textureComponent.ID]
			if !ok {
				continue
			}

			meshComponent, err := ecs.GetComponent[meshcomponent.Mesh](world, entity)
			if err != nil {
				continue
			}
			meshAsset, ok := m.meshes[meshComponent.ID]
			if !ok {
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

			textureAsset.Use()

			meshAsset.Draw()
		}
		return nil
	}

	// projectionsEntities := map[projection.UsedProjection][]ecs.EntityID{}
	//
	// for _, entity := range entities {
	// 	usedProjection, err := ecs.GetComponent[projection.UsedProjection](world, entity)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	arr, _ := projectionsEntities[usedProjection]
	// 	arr = append(arr, entity)
	// 	projectionsEntities[usedProjection] = arr
	// }
	//
	// for usedProjection, entities := range projectionsEntities {
	// 	query := world.QueryEntitiesWithComponents(usedProjection.ProjectionComponent)
	// 	cameraEntities := query.Entities()
	// 	if len(cameraEntities) != 1 {
	// 		return projection.ErrWorldShouldHaveOneProjection
	// 	}
	// 	camera := cameraEntities[0]
	//
	// 	projectionComponent, err := usedProjection.GetCameraProjection(world, camera)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	cameraTransformComponent, err := ecs.GetComponent[transform.Transform](world, camera)
	// 	if err != nil {
	// 		return errors.Join(errors.New("camera misses transform component"), err)
	// 	}
	//
	// 	projectionMat4 := projectionComponent.Mat4()
	// 	cameraTransformMat4 := projectionComponent.ViewMat4(cameraTransformComponent)
	//
	// 	mvp := projectionMat4.Mul4(cameraTransformMat4)
	//
	// 	for _, entity := range entities {
	// 		transformComponent, err := ecs.GetComponent[transform.Transform](world, entity)
	// 		if err != nil {
	// 			continue
	// 		}
	//
	// 		textureComponent, err := ecs.GetComponent[texturecomponent.Texture](world, entity)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		textureAsset, ok := m.textures[textureComponent.ID]
	// 		if !ok {
	// 			continue
	// 		}
	//
	// 		meshComponent, err := ecs.GetComponent[meshcomponent.Mesh](world, entity)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		meshAsset, ok := m.meshes[meshComponent.ID]
	// 		if !ok {
	// 			continue
	// 		}
	//
	// 		// locations
	// 		locations, err := program.GetProgramLocations[locations](p)
	// 		if err != nil {
	// 			continue
	// 		}
	//
	// 		mvp := mvp.Mul4(transformComponent.Mat4())
	// 		gl.UniformMatrix4fv(locations.Mvp, 1, false, &mvp[0])
	// 		textureAsset.Use()
	// 		meshAsset.Draw()
	// 	}
	// }

	return nil
}
