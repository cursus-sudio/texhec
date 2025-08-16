package render

import (
	"frontend/engine/components/material"
	"frontend/engine/components/mesh"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/frames"
)

type RenderSystem struct {
	world  ecs.World
	assets assets.Assets

	renderables map[assets.AssetID]*renderable
}

func NewRenderSystem(
	world ecs.World,
	assetsService assets.Assets,
) RenderSystem {
	renderedEntities := map[ecs.EntityID]assets.AssetID{}
	renderables := map[assets.AssetID]*renderable{}

	onAdd := func(entities []ecs.EntityID) {
		for _, entity := range entities {
			materialComponent, err := ecs.GetComponent[material.Material](world, entity)
			if err != nil {
				continue
			}
			materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](assetsService, materialComponent.ID)
			if err != nil {
				continue
			}

			renderedEntities[entity] = materialComponent.ID
			if existing, ok := renderables[materialComponent.ID]; ok {
				existing.Entities = append(existing.Entities, entity)
			} else {
				renderables[materialComponent.ID] = &renderable{
					materialAsset,
					[]ecs.EntityID{entity},
				}
			}
		}
	}
	onRemove := func(entities []ecs.EntityID) {
		for _, entity := range entities {
			assetID, ok := renderedEntities[entity]
			if !ok {
				continue
			}
			renderable, ok := renderables[assetID]
			if !ok {
				continue
			}
			newEntities := make([]ecs.EntityID, 0)
			for _, renderedEntity := range renderable.Entities {
				if renderedEntity != entity {
					newEntities = append(newEntities, renderedEntity)
				}
			}
		}
	}

	liveQuery := world.GetEntitiesWithComponentsQuery(
		ecs.GetComponentType(mesh.Mesh{}),
		ecs.GetComponentType(material.Material{}),
	)
	liveQuery.OnAdd(onAdd)
	liveQuery.OnRemove(onRemove)

	return RenderSystem{
		world:  world,
		assets: assetsService,

		renderables: renderables,
	}
}

type renderable struct {
	Material material.MaterialCachedAsset
	Entities []ecs.EntityID
}

func (s *RenderSystem) Listen(args frames.FrameEvent) error {
	for _, material := range s.renderables {
		if err := material.Material.Render(s.world, material.Entities); err != nil {
			return err
		}
	}

	return nil
}
