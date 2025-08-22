package texturematerial

import (
	"errors"
	meshcomponent "frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/datastructures"
	"frontend/services/ecs"
	"frontend/services/graphics/buffers"
	"frontend/services/graphics/program"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"image"
	"sync"

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

	mutex *sync.RWMutex

	entities        datastructures.Set[ecs.EntityID]
	modelBuffer     buffers.Buffer[mgl32.Mat4]
	modelProjBuffer buffers.Buffer[int32]
	modelTexBuffer  buffers.Buffer[int32]
	cmdBuffer       buffers.Buffer[DrawElementsIndirectCommand]
	// currently there is 1 entity 1 command
	// TODO add instancing

	projBuffer buffers.Buffer[mgl32.Mat4]

	// for modelCmdBuffer and modelTexBuffer
	meshes   map[assets.AssetID]int32
	textures map[assets.AssetID]int32

	// render data
	mesh    vao.VAO
	texture uint32

	packedMesh []MeshRange

	projections map[ecs.ComponentType]int32
}

type textureMaterialServices struct {
	window        window.Api
	assetsStorage assets.AssetsStorage
	locations     locations
	console       console.Console

	*renderCache
}

func (m *textureMaterialServices) cleanUp() {
	if m.cachedForWorld == nil {
		return
	}
	//

	m.mutex = &sync.RWMutex{}

	m.cachedForWorld = nil
	m.mesh.Release()
	m.mesh = nil
	m.packedMesh = nil
	m.meshes = map[assets.AssetID]int32{}

	m.modelBuffer.Release()
	m.modelProjBuffer.Release()
	m.modelTexBuffer.Release()
	m.cmdBuffer.Release()
	m.projBuffer.Release()
	m.entities = nil
	m.modelBuffer = nil
	m.modelProjBuffer = nil
	m.modelTexBuffer = nil
	m.cmdBuffer = nil
	m.projBuffer = nil

	gl.DeleteTextures(1, &m.texture)
	m.texture = 0
	m.textures = map[assets.AssetID]int32{}

	m.projections = map[ecs.ComponentType]int32{}
}

func (m *textureMaterialServices) init(world ecs.World, p program.Program) error {
	m.cachedForWorld = world
	m.mesh = nil
	m.packedMesh = nil
	m.meshes = map[assets.AssetID]int32{}

	m.mutex = &sync.RWMutex{}
	m.texture = 0
	m.textures = map[assets.AssetID]int32{}

	m.projections = map[ecs.ComponentType]int32{
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
				m.textures[assetID] = int32(i)
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
				m.meshes[assetID] = int32(i)
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

		m.entities = datastructures.NewSet[ecs.EntityID]()

		gl.GenBuffers(1, &buffer)
		m.cmdBuffer = buffers.NewBuffer[DrawElementsIndirectCommand](
			gl.DRAW_INDIRECT_BUFFER, gl.DYNAMIC_DRAW, buffer)

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, buffer)
		m.modelTexBuffer = buffers.NewBuffer[int32](
			gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, buffer)
		m.modelBuffer = buffers.NewBuffer[mgl32.Mat4](
			gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, buffer)
		m.modelProjBuffer = buffers.NewBuffer[int32](
			gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

		// proj buffer

		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, buffer)
		m.projBuffer = buffers.NewBuffer[mgl32.Mat4](
			gl.SHADER_STORAGE_BUFFER, gl.STATIC_DRAW, buffer)
	}
	for i := 0; i < 3; i++ {
		m.projBuffer.Add(mgl32.Ident4())
	}

	for projectionType, projectionIndex := range m.projections {
		query := world.QueryEntitiesWithComponents(
			projectionType,
			ecs.GetComponentType(transform.Transform{}),
		)

		onChange := func(_ []ecs.EntityID) {
			m.mutex.Lock()
			defer m.mutex.Unlock()
			entities := query.Entities()
			if len(entities) != 1 {
				return // projection.ErrWorldShouldHaveOneProjection
			}
			camera := entities[0]

			anyProj, err := world.GetComponent(camera, projectionType)
			if err != nil {
				return // err
			}

			projectionComponent, ok := anyProj.(projection.Projection)
			if !ok {
				return // projection.ErrExpectedUsedProjectionToImplementProjection
			}

			cameraTransformComponent, err := ecs.GetComponent[transform.Transform](world, camera)
			if err != nil {
				return // errors.Join(errors.New("camera misses transform component"), err)
			}

			projectionMat4 := projectionComponent.Mat4()
			cameraTransformMat4 := projectionComponent.ViewMat4(cameraTransformComponent)

			mvp := projectionMat4.Mul4(cameraTransformMat4)
			m.projBuffer.Set(int(projectionIndex), mvp)
		}
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
			ecs.GetComponentType(projection.Visible{}),
		)
		onChange := func(entities []ecs.EntityID) {
			m.mutex.Lock()
			defer m.mutex.Unlock()
			for _, entity := range entities {
				transformComponent, err := ecs.GetComponent[transform.Transform](world, entity)
				if err != nil {
					continue
				}
				model := transformComponent.Mat4()

				textureComponent, err := ecs.GetComponent[texturecomponent.Texture](world, entity)
				if err != nil {
					continue
				}
				textureIndex, ok := m.textures[textureComponent.ID]
				if !ok {
					continue
				}

				meshComponent, err := ecs.GetComponent[meshcomponent.Mesh](world, entity)
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
					cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(len(m.entities.Get())))
					m.entities.Add(entity)
					m.cmdBuffer.Add(cmd)
					m.modelTexBuffer.Add(textureIndex)
					m.modelBuffer.Add(model)
					m.modelProjBuffer.Add(projectionIndex)
					continue
				}
				cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(index))
				m.cmdBuffer.Set(index, cmd)
				m.modelTexBuffer.Set(index, textureIndex)
				m.modelBuffer.Set(index, model)
				m.modelProjBuffer.Set(index, projectionIndex)
			}
		}
		onRemove := func(entities []ecs.EntityID) {
			m.mutex.Lock()
			defer m.mutex.Unlock()
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

	{
		m.mutex.Lock()

		m.cmdBuffer.Flush()
		m.modelTexBuffer.Flush()
		m.modelBuffer.Flush()
		m.modelProjBuffer.Flush()
		m.projBuffer.Flush()

		m.mutex.Unlock()
	}

	p.Use()
	m.mesh.Use()
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, m.texture)
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, m.cmdBuffer.ID())
	gl.MultiDrawElementsIndirect(gl.TRIANGLES, gl.UNSIGNED_INT, nil, int32(len(m.cmdBuffer.Data())), 0)

	return nil
}
