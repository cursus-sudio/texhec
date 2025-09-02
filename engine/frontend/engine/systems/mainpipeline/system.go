package mainpipeline

import (
	"errors"
	"fmt"
	meshcomponent "frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/systems/render"
	"frontend/engine/tools/worldmesh"
	"frontend/engine/tools/worldtexture"
	"frontend/services/assets"
	"frontend/services/datastructures"
	"frontend/services/ecs"
	"frontend/services/graphics"
	"frontend/services/graphics/buffers"
	"frontend/services/media/window"
	"shared/services/logger"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

type PipelineComponent struct{}

//

type System struct {
	world         ecs.World
	window        window.Api
	assetsStorage assets.AssetsStorage
	logger        logger.Logger

	entitiesQueryAdditionalArguments []ecs.ComponentType
}

func NewSystem(
	world ecs.World,
	window window.Api,
	assetsStorage assets.AssetsStorage,
	logger logger.Logger,
	entitiesQueryAdditionalArguments []ecs.ComponentType,
) (*System, error) {
	system := &System{
		world: world,

		window:        window,
		assetsStorage: assetsStorage,
		logger:        logger,

		entitiesQueryAdditionalArguments: entitiesQueryAdditionalArguments,
	}

	projections := datastructures.NewSet[ecs.ComponentType]()
	projections.Add(ecs.GetComponentType(projection.Ortho{}), ecs.GetComponentType(projection.Perspective{}))

	register, err := newRegister(projections)
	if err != nil {
		return nil, err
	}
	world.SaveRegister(register)
	system.modifyRegisterOnChanges()
	return system, nil
}

func listenToProjectionChanges[Projection projection.Projection](
	m *System,
	mutex *sync.RWMutex,
	buffer buffers.Buffer[mgl32.Mat4],
) {

	var zero Projection
	query := m.world.QueryEntitiesWithComponents(
		ecs.GetComponentType(zero),
		ecs.GetComponentType(transform.Transform{}),
	)
	onChange := func(_ []ecs.EntityID) {
		mutex.Lock()
		defer mutex.Unlock()
		entities := query.Entities()
		for len(buffer.Get()) < len(entities) {
			buffer.Add(mgl32.Mat4{})
		}
		for len(buffer.Get()) > len(entities) {
			buffer.Remove(len(buffer.Get()) - 1)
		}
		for i, camera := range entities {
			projectionComponent, err := ecs.GetComponent[Projection](m.world, camera)
			if err != nil {
				m.logger.Error(err)
				return
			}

			cameraTransformComponent, err := ecs.GetComponent[transform.Transform](m.world, camera)
			if err != nil {
				m.logger.Error(errors.New("camera misses transform component"))
				return
			}

			projectionMat4 := projectionComponent.Mat4()
			cameraTransformMat4 := projectionComponent.ViewMat4(cameraTransformComponent)

			mvp := projectionMat4.Mul4(cameraTransformMat4)
			buffer.Set(i, mvp)
		}
	}
	query.OnAdd(onChange)
	query.OnChange(onChange)
	query.OnRemove(onChange)
}

func (m *System) modifyRegisterOnChanges() error {
	register, err := ecs.GetRegister[pipelineRegister](m.world)
	if err != nil {
		return err
	}
	listenToProjectionChanges[projection.Ortho](m, register.mutex, register.buffers.orthoBuffer)
	listenToProjectionChanges[projection.Perspective](m, register.mutex, register.buffers.perspectiveBuffer)

	// change entities buffer
	{
		onChange := func(entities []ecs.EntityID) {
			register.mutex.Lock()
			defer register.mutex.Unlock()

			mesh, err := ecs.GetRegister[worldmesh.WorldMeshRegister[Vertex]](m.world)
			if err != nil {
				m.logger.Error(err)
				return
			}

			texture, err := ecs.GetRegister[worldtexture.WorldTextureRegister](m.world)
			if err != nil {
				m.logger.Error(err)
				return
			}

			for _, entity := range entities {
				transformComponent, err := ecs.GetComponent[transform.Transform](m.world, entity)
				if err != nil {
					continue
				}
				model := transformComponent.Mat4()

				textureComponent, err := ecs.GetComponent[texturecomponent.Texture](m.world, entity)
				if err != nil {
					continue
				}
				textureIndex, ok := texture.Assets.GetIndex(textureComponent.ID)
				if !ok {
					m.logger.Error(fmt.Errorf(
						"pipeline cannot render entity with texture which isn't in world texture",
					))
					continue
				}

				meshComponent, err := ecs.GetComponent[meshcomponent.Mesh](m.world, entity)
				if err != nil {
					continue
				}
				meshRange, ok := mesh.Ranges[meshComponent.ID]
				if !ok {
					m.logger.Error(fmt.Errorf(
						"pipeline cannot render entity with mesh which isn't in world mesh",
					))
					continue
				}

				usedProjection, err := ecs.GetComponent[projection.UsedProjection](m.world, entity)
				if err != nil {
					continue
				}

				projectionIndex, ok := register.projections.GetIndex(usedProjection.ProjectionComponent)
				if !ok {
					m.logger.Error(fmt.Errorf(
						"pipeline doesn't handle \"%s\" projection",
						usedProjection.ProjectionComponent.String(),
					))
					continue
				}

				cmd := graphics.NewDrawElementsIndirectCommand(
					meshRange.IndexCount,
					1,
					meshRange.FirstIndex,
					meshRange.FirstVertex,
					0,
				)
				register.buffers.Upsert(
					entity,
					cmd,
					int32(textureIndex),
					model,
					int32(projectionIndex),
				)
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(len(m.entities.Get())))
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(index))
			}
		}

		onRemove := func(entities []ecs.EntityID) {
			register.mutex.Lock()
			defer register.mutex.Unlock()
			register.buffers.Remove(entities)
		}

		query := m.world.QueryEntitiesWithComponents(
			append(
				m.entitiesQueryAdditionalArguments,
				ecs.GetComponentType(PipelineComponent{}),
				ecs.GetComponentType(transform.Transform{}),
				ecs.GetComponentType(projection.UsedProjection{}),
				ecs.GetComponentType(meshcomponent.Mesh{}),
				ecs.GetComponentType(texturecomponent.Texture{}),
			)...,
		)

		query.OnAdd(onChange)
		query.OnChange(onChange)
		query.OnRemove(onRemove)
	}
	return nil
}

var changes0 int = 0

func (m *System) Listen(render.RenderEvent) error {
	pipelineRegister, err := ecs.GetRegister[pipelineRegister](m.world)
	if err != nil {
		return err
	}
	mesh, err := ecs.GetRegister[worldmesh.WorldMeshRegister[Vertex]](m.world)
	if err != nil {
		return err
	}
	texture, err := ecs.GetRegister[worldtexture.WorldTextureRegister](m.world)
	if err != nil {
		return err
	}

	pipelineRegister.mutex.Lock()
	pipelineRegister.buffers.Flush()
	pipelineRegister.mutex.Unlock()

	pipelineRegister.program.Use()
	mesh.Mesh.Use()
	texture.Use()
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, pipelineRegister.buffers.cmdBuffer.ID())
	cmdBufferLen := len(pipelineRegister.buffers.cmdBuffer.Get())
	gl.MultiDrawElementsIndirect(
		gl.TRIANGLES,
		gl.UNSIGNED_INT,
		nil,
		int32(cmdBufferLen),
		0,
	)

	return nil
}
