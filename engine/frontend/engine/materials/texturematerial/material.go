package texturematerial

import (
	"errors"
	"frontend/engine/components/mesh"
	meshcomponent "frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/materials/texturematerial/arrays"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

//

type locations struct {
	// Mvp int32 `uniform:"mvp"`
}

type renderCache struct {
	cachedForWorld ecs.World
	// entitiesLiveQuery ecs.LiveQuery

	//

	entities        arrays.IndexTracker[ecs.EntityID]
	modelBuffer     arrays.Buffer[mgl32.Mat4]
	modelProjBuffer arrays.Buffer[int]
	modelTexBuffer  arrays.Buffer[int]
	cmdBuffer       arrays.Buffer[DrawElementsIndirectCommand]
	// currently there is 1 entity 1 command
	// TODO add instancing

	projBuffer arrays.Buffer[mgl32.Mat4]

	// for modelCmdBuffer and modelTexBuffer
	meshes   map[assets.AssetID]int
	textures map[assets.AssetID]int

	// render data
	mesh    vao.VAO
	texture uint32

	packedMesh []MeshRange

	projections map[ecs.ComponentType]int
}

type textureMaterialServices struct {
	window        window.Api
	assetsStorage assets.AssetsStorage
	locations     locations

	*renderCache
}

func (m *textureMaterialServices) cleanUp() {
	if m.cachedForWorld == nil {
		return
	}
	m.cachedForWorld = nil
	m.mesh.Release()
	m.mesh = nil
	m.packedMesh = nil
	m.meshes = map[assets.AssetID]int{}

	gl.DeleteTextures(1, &m.texture)
	m.texture = 0
	m.textures = map[assets.AssetID]int{}
}

func (m *textureMaterialServices) init(world ecs.World, p program.Program) error {
	m.cachedForWorld = world
	m.mesh = nil
	m.packedMesh = nil
	m.meshes = map[assets.AssetID]int{}

	m.texture = 0
	m.textures = map[assets.AssetID]int{}

	m.projections = map[ecs.ComponentType]int{
		ecs.GetComponentType(projection.Ortho{}):       0,
		ecs.GetComponentType(projection.Perspective{}): 1,
	}

	{ // generate texture and mesh
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

		{
			images := make([]image.Image, len(worldComponent.Textures))
			w, h := 0, 0
			for i, assetID := range worldComponent.Textures {
				textureAsset, err := assets.StorageGet[texturecomponent.TextureStorageAsset](m.assetsStorage, assetID)
				if err != nil {
					continue
				}

				image := textureAsset.Image()
				if i == 0 {
					w, h = image.Bounds().Dx(), image.Bounds().Dy()
				}
				if w != image.Bounds().Dx() || h != image.Bounds().Dy() {
					return ErrTexturesHaveToShareSize
				}
				images[i] = image
				m.textures[assetID] = i
			}
			m.texture = CreateTexs(w, h, images)
			texLoc := gl.GetUniformLocation(p.ID(), gl.Str("texs\x00"))
			gl.Uniform1i(texLoc, 1)
		}

		{
			meshes := make([]Mesh, len(worldComponent.Meshes))
			for i, assetID := range worldComponent.Meshes {
				meshAsset, err := assets.StorageGet[meshcomponent.MeshStorageAsset](m.assetsStorage, assetID)
				if err != nil {
					continue
				}
				mesh := Mesh{
					meshAsset.Verticies(),
					meshAsset.Indicies(),
				}
				meshes[i] = mesh
				m.meshes[assetID] = i
			}
			packedMesh := Pack(meshes...)

			vbo := vbo.NewVBO()
			vbo.SetVertices(packedMesh.vertices)

			ebo := ebo.NewEBO()
			ebo.SetIndices(packedMesh.indices)

			vao := vao.NewVAO(vbo, ebo)
			m.mesh = vao
			m.packedMesh = packedMesh.ranges
		}
	}

	{
		var buffer uint32

		// model buffer

		m.entities = arrays.NewIndexTracker(
			arrays.NewArray[ecs.EntityID](),
		)

		gl.GenBuffers(1, &buffer)
		m.cmdBuffer = arrays.NewBuffer[DrawElementsIndirectCommand](
			gl.DRAW_INDIRECT_BUFFER, gl.DYNAMIC_DRAW, buffer)

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, buffer)
		m.modelTexBuffer = arrays.NewBuffer[int](
			gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, buffer)
		m.modelBuffer = arrays.NewBuffer[mgl32.Mat4](
			gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, buffer)
		m.modelProjBuffer = arrays.NewBuffer[int](
			gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

		// proj buffer

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, buffer)
		m.projBuffer = arrays.NewBuffer[mgl32.Mat4](
			gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)
	}

	for projectionType, projectionIndex := range m.projections {
		m.projBuffer.Add(mgl32.Ident4())

		onChange := func(entities []ecs.EntityID) {
			if len(entities) != 1 {
				return // projection.ErrWorldShouldHaveOneProjection
			}
			camera := entities[0]

			usedProjection := projection.UsedProjection{ProjectionComponent: projectionType}
			projectionComponent, err := usedProjection.GetCameraProjection(world, camera)
			if err != nil {
				return // err
			}

			cameraTransformComponent, err := ecs.GetComponent[transform.Transform](world, camera)
			if err != nil {
				return // errors.Join(errors.New("camera misses transform component"), err)
			}

			projectionMat4 := projectionComponent.Mat4()
			cameraTransformMat4 := projectionComponent.ViewMat4(cameraTransformComponent)

			mvp := projectionMat4.Mul4(cameraTransformMat4)
			m.projBuffer.Set(projectionIndex, mvp)
		}
		query := world.QueryEntitiesWithComponents(projectionType)
		query.OnAdd(onChange)
		query.OnChange(onChange)
	}

	{
		query := world.QueryEntitiesWithComponents(
			ecs.GetComponentType(TextureMaterialComponent{}),
			ecs.GetComponentType(meshcomponent.Mesh{}),
			ecs.GetComponentType(projection.UsedProjection{}),
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(texturecomponent.Texture{}),
		)
		onChange := func(entities []ecs.EntityID) {
			for i, entity := range entities {
				transformComponent, err := ecs.GetComponent[transform.Transform](world, entity)
				if err != nil {
					continue
				}
				model := transformComponent.Mat4()

				textureComponent, err := ecs.GetComponent[texture.Texture](world, entity)
				if err != nil {
					continue
				}
				textureIndex, ok := m.textures[textureComponent.ID]
				if !ok {
					continue
				}

				meshComponent, err := ecs.GetComponent[mesh.Mesh](world, entity)
				if err != nil {
					continue
				}
				meshIndex, ok := m.meshes[meshComponent.ID]
				if !ok {
					continue
				}
				meshRange := m.packedMesh[meshIndex]
				if !ok {
					continue
				}
				cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(i))

				usedProjection, err := ecs.GetComponent[projection.UsedProjection](world, entity)
				if err != nil {
					continue
				}
				projectionIndex, ok := m.projections[usedProjection.ProjectionComponent]
				if !ok {
					continue
				}

				index, ok := m.entities.GetIndex(entity)
				if !ok {
					m.entities.Add(entity)
					m.cmdBuffer.Add(cmd)
					m.modelTexBuffer.Add(textureIndex)
					m.modelBuffer.Add(model)
					m.modelProjBuffer.Add(projectionIndex)
					continue
				}
				m.cmdBuffer.Set(index, cmd)
				m.modelTexBuffer.Set(index, textureIndex)
				m.modelBuffer.Set(index, model)
				m.modelProjBuffer.Set(index, projectionIndex)
			}
		}
		onRemove := func(entities []ecs.EntityID) {
			for _, entity := range entities {
				index, ok := m.entities.GetIndex(entity)
				if !ok {
					continue
				}
				m.entities.Remove(index)
				m.cmdBuffer.Remove(index)
				m.modelTexBuffer.Remove(index)
				m.modelBuffer.Remove(index)
				m.modelProjBuffer.Remove(index)
			}
		}
		query.OnAdd(onChange)
		query.OnChange(onChange)
		query.OnRemove(onRemove)
	}

	return nil
}

func (m *textureMaterialServices) render(world ecs.World, p program.Program) error {
	reCache := m.cachedForWorld != world
	if reCache {
		m.cleanUp()
		if err := m.init(world, p); err != nil {
			return err
		}
	}

	m.cmdBuffer.Flush()
	m.modelTexBuffer.Flush()
	m.modelBuffer.Flush()
	m.modelProjBuffer.Flush()
	m.projBuffer.Flush()

	p.Use()
	m.mesh.Use()
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, m.texture)
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, m.cmdBuffer.ID())
	gl.MultiDrawElementsIndirect(gl.TRIANGLES, gl.UNSIGNED_INT, nil, int32(len(m.cmdBuffer.Data())), 0)

	return nil
}
