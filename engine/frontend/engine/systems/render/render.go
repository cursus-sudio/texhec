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
}

func NewRenderSystem(
	world ecs.World,
	assets assets.Assets,
) RenderSystem {
	return RenderSystem{
		world:  world,
		assets: assets,
	}
}

type renderable struct {
	Material material.MaterialCachedAsset
	Entities []ecs.EntityID
}

func (s *RenderSystem) Listen(args frames.FrameEvent) error {
	renderables := map[assets.AssetID]renderable{}

	renderableEntities := s.world.GetEntitiesWithComponents(
		ecs.GetComponentType(mesh.Mesh{}),
		ecs.GetComponentType(material.Material{}),
	)

	for _, entity := range renderableEntities {
		materialComponent, err := ecs.GetComponent[material.Material](s.world, entity)
		if err != nil {
			continue
		}
		materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](s.assets, materialComponent.ID)
		if err != nil {
			return err
		}

		if existing, ok := renderables[materialComponent.ID]; ok {
			existing.Entities = append(existing.Entities, entity)
			renderables[materialComponent.ID] = existing
		} else {
			renderables[materialComponent.ID] = renderable{
				materialAsset,
				[]ecs.EntityID{entity},
			}
		}
	}

	for _, material := range renderables {
		if err := material.Material.Render(s.world, material.Entities); err != nil {
			return err
		}
	}

	return nil
}
