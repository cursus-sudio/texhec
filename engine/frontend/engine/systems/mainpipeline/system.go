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
	"frontend/services/ecs"
	"frontend/services/graphics"
	"frontend/services/media/window"
	"shared/services/logger"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

// this component says which entities use material
type MaterialComponent struct{}

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

	register, err := newMaterialWorldRegistry(
		map[ecs.ComponentType]int32{
			ecs.GetComponentType(projection.Ortho{}):       0,
			ecs.GetComponentType(projection.Perspective{}): 1,
		},
	)
	if err != nil {
		return nil, err
	}
	world.SaveRegister(register)
	system.modifyRegisterOnChanges()
	return system, nil
}

func (m *System) modifyRegisterOnChanges() error {
	register, err := ecs.GetRegister[materialWorldRegister](m.world)
	if err != nil {
		return err
	}
	// modify projections buffer
	for projectionType, projectionIndex := range register.projections {
		query := m.world.QueryEntitiesWithComponents(
			projectionType,
			ecs.GetComponentType(transform.Transform{}),
		)
		onChange := func(_ []ecs.EntityID) {
			register.mutex.Lock()
			defer register.mutex.Unlock()
			entities := query.Entities()
			if len(entities) != 1 {
				m.logger.Error(projection.ErrWorldShouldHaveOneProjection)
				return
			}
			camera := entities[0]

			anyProj, err := m.world.GetComponent(camera, projectionType)
			if err != nil {
				m.logger.Error(err)
				return
			}

			projectionComponent, ok := anyProj.(projection.Projection)
			if !ok {
				m.logger.Error(projection.ErrExpectedUsedProjectionToImplementProjection)
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
			register.buffers.projBuffer.Set(int(projectionIndex), mvp)
		}
		query.OnAdd(onChange)
		query.OnChange(onChange)
	}

	//

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
				// textureIndex, ok := register.textures[textureComponent.ID]
				textureIndex, ok := texture.Assets.GetIndex(textureComponent.ID)
				if !ok {
					m.logger.Error(fmt.Errorf(
						"material cannot render entity with texture which isn't in WorldTextureMaterialComponent",
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
						"material cannot render entity with mesh which isn't in WorldTextureMaterialComponent",
					))
					continue
				}

				usedProjection, err := ecs.GetComponent[projection.UsedProjection](m.world, entity)
				if err != nil {
					continue
				}
				projectionIndex, ok := register.projections[usedProjection.ProjectionComponent]
				if !ok {
					m.logger.Error(fmt.Errorf(
						"material doesn't handle \"%s\" projection",
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
					projectionIndex,
				)
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(len(m.entities.Get())))
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(index))
			}
		}

		query := m.world.QueryEntitiesWithComponents(
			append(
				m.entitiesQueryAdditionalArguments,
				ecs.GetComponentType(MaterialComponent{}),
				ecs.GetComponentType(transform.Transform{}),
				ecs.GetComponentType(projection.UsedProjection{}),
				ecs.GetComponentType(meshcomponent.Mesh{}),
				ecs.GetComponentType(texturecomponent.Texture{}),
			)...,
		)

		query.OnAdd(onChange)
		query.OnChange(onChange)
		query.OnRemove(register.buffers.Remove)
	}
	return nil
}

func (m *System) Listen(render.RenderEvent) error {
	materialRegister, err := ecs.GetRegister[materialWorldRegister](m.world)
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

	materialRegister.mutex.Lock()
	materialRegister.buffers.Flush()
	materialRegister.mutex.Unlock()

	materialRegister.program.Use()
	mesh.Mesh.Use()
	texture.Bind()
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, materialRegister.buffers.cmdBuffer.ID())
	gl.MultiDrawElementsIndirect(
		gl.TRIANGLES,
		gl.UNSIGNED_INT,
		nil,
		int32(len(materialRegister.buffers.cmdBuffer.Get())),
		0,
	)

	return nil
}
