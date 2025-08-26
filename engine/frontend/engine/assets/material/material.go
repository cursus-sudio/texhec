package material

import (
	"errors"
	"fmt"
	meshcomponent "frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/buffers"
	"frontend/services/graphics/program"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/media/window"
	"image"
	"shared/services/logger"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

type materialWorldCache struct {
	world ecs.World
	mutex *sync.RWMutex

	entitiesBuffers *entitiesBuffers

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

type materialCache struct {
	window        window.Api
	assetsStorage assets.AssetsStorage
	logger        logger.Logger

	entitiesQueryAdditionalArguments []ecs.ComponentType

	*materialWorldCache
}

func (m *materialCache) cleanUp() {
	if m.world == nil {
		return
	}

	m.mutex = &sync.RWMutex{}
	m.world = nil

	m.mesh.Release()
	m.mesh = nil
	m.packedMesh = nil
	m.meshes = map[assets.AssetID]int32{}

	m.entitiesBuffers.CleanUp()
	m.entitiesBuffers = nil

	m.projBuffer.Release()
	m.projBuffer = nil

	gl.DeleteTextures(1, &m.texture)
	m.texture = 0
	m.textures = map[assets.AssetID]int32{}

	m.projections = map[ecs.ComponentType]int32{}
}

func (m *materialCache) init(world ecs.World, p program.Program) error {
	m.world = world
	m.mesh = nil
	m.packedMesh = nil
	m.meshes = map[assets.AssetID]int32{}

	m.mutex = &sync.RWMutex{}
	query := world.QueryEntitiesWithComponents(
		append(
			m.entitiesQueryAdditionalArguments,
			ecs.GetComponentType(TextureMaterialComponent{}),
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(projection.UsedProjection{}),
			ecs.GetComponentType(meshcomponent.Mesh{}),
			ecs.GetComponentType(texturecomponent.Texture{}),
		)...,
	)
	m.world = world

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
			meshes := make([]Mesh[Vertex], len(worldComponent.Meshes))
			for i, assetID := range worldComponent.Meshes {
				meshAsset, err := assets.StorageGet[meshcomponent.MeshStorageAsset[Vertex]](m.assetsStorage, assetID)
				if err != nil {
					continue
				}
				mesh := Mesh[Vertex]{
					meshAsset.Verticies(),
					meshAsset.Indicies(),
				}
				meshes[i] = mesh
				m.meshes[assetID] = int32(i)
			}
			packedMesh := Pack(meshes...)

			vbo := NewVBO()
			vbo.SetVertices(packedMesh.vertices)

			ebo := ebo.NewEBO()
			ebo.SetIndices(packedMesh.indices)

			vao := vao.NewVAO(vbo, ebo)
			m.mesh = vao
			m.packedMesh = packedMesh.ranges
		}
	}

	{
		m.entitiesBuffers = newEntitiesBuffers()
		var buffer uint32
		gl.GenBuffers(1, &buffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, buffer)
		m.projBuffer = buffers.NewBuffer[mgl32.Mat4](gl.SHADER_STORAGE_BUFFER, gl.STATIC_DRAW, buffer)
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
				m.logger.Error(projection.ErrWorldShouldHaveOneProjection)
				return
			}
			camera := entities[0]

			anyProj, err := world.GetComponent(camera, projectionType)
			if err != nil {
				m.logger.Error(err)
				return
			}

			projectionComponent, ok := anyProj.(projection.Projection)
			if !ok {
				m.logger.Error(projection.ErrExpectedUsedProjectionToImplementProjection)
				return
			}

			cameraTransformComponent, err := ecs.GetComponent[transform.Transform](world, camera)
			if err != nil {
				m.logger.Error(errors.New("camera misses transform component"))
				return
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
					m.logger.Error(fmt.Errorf(
						"material cannot render entity with texture which isn't in WorldTextureMaterialComponent",
					))
					continue
				}

				meshComponent, err := ecs.GetComponent[meshcomponent.Mesh](world, entity)
				if err != nil {
					continue
				}
				meshIndex, ok := m.meshes[meshComponent.ID]
				if !ok {
					m.logger.Error(fmt.Errorf(
						"material cannot render entity with mesh which isn't in WorldTextureMaterialComponent",
					))
					continue
				}
				meshRange := m.packedMesh[meshIndex]

				usedProjection, err := ecs.GetComponent[projection.UsedProjection](world, entity)
				if err != nil {
					continue
				}
				projectionIndex, ok := m.projections[usedProjection.ProjectionComponent]
				if !ok {
					m.logger.Error(fmt.Errorf(
						"material doesn't handle \"%s\" projection",
						usedProjection.ProjectionComponent.String(),
					))
					continue
				}

				cmd := meshRange.DrawCommand(1, 0)
				m.entitiesBuffers.Upsert(
					entity,
					cmd,
					textureIndex,
					model,
					projectionIndex,
				)
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(len(m.entities.Get())))
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(index))
			}
		}

		query.OnAdd(onChange)
		query.OnChange(onChange)
		query.OnRemove(m.entitiesBuffers.Remove)
	}

	return nil
}

func (m *materialCache) render(world ecs.World, p program.Program) error {
	reCache := m.world != world
	if reCache {
		m.cleanUp()
		if err := m.init(world, p); err != nil {
			return err
		}
	}

	{
		m.mutex.Lock()

		m.entitiesBuffers.Flush()
		m.projBuffer.Flush()

		m.mutex.Unlock()
	}

	p.Use()
	m.mesh.Use()
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, m.texture)
	cmds := m.entitiesBuffers.cmdBuffer
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, cmds.ID())
	gl.MultiDrawElementsIndirect(gl.TRIANGLES, gl.UNSIGNED_INT, nil, int32(len(cmds.Get())), 0)

	return nil
}
