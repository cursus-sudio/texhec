package material

import (
	"errors"
	"fmt"
	meshcomponent "frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/tools/worldmesh"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics"
	"frontend/services/graphics/program"
	"frontend/services/media/window"
	"shared/services/logger"
)

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

type materialCache struct {
	window        window.Api
	assetsStorage assets.AssetsStorage
	logger        logger.Logger

	entitiesQueryAdditionalArguments []ecs.ComponentType
}

func (m *materialCache) modifyRegisterOnChanges(
	world ecs.World,
	register materialWorldRegister,
) {
	// modify projections buffer
	for projectionType, projectionIndex := range register.projections {
		query := world.QueryEntitiesWithComponents(
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

			geometry, err := ecs.GetRegister[worldmesh.WorldMeshRegister[Vertex]](world)
			if err != nil {
				m.logger.Error(err)
				return
			}

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
				textureIndex, ok := register.textures[textureComponent.ID]
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
				meshRange, ok := geometry.Ranges[meshComponent.ID]
				if !ok {
					m.logger.Error(fmt.Errorf(
						"material cannot render entity with mesh which isn't in WorldTextureMaterialComponent",
					))
					continue
				}

				usedProjection, err := ecs.GetComponent[projection.UsedProjection](world, entity)
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
					textureIndex,
					model,
					projectionIndex,
				)
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(len(m.entities.Get())))
				// cmd := NewDrawElementsIndirectCommand(meshRange, 1, uint32(index))
			}
		}

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

		query.OnAdd(onChange)
		query.OnChange(onChange)
		query.OnRemove(register.buffers.Remove)
	}
}

func (m *materialCache) render(world ecs.World, p program.Program) error {
	register, err := ecs.GetRegister[materialWorldRegister](world)
	if err != nil {
		register, err = createRegister(
			world,
			p,
			m.assetsStorage,
		)
		if err != nil {
			return err
		}
		m.modifyRegisterOnChanges(world, register)
		world.SaveRegister(register)
	}

	return register.Render(world, p)
}
