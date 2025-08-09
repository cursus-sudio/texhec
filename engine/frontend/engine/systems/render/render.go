package render

import (
	"frontend/engine/components/material"
	"frontend/engine/components/mesh"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/frames"
)

type RenderSystem struct {
	World  ecs.World
	Assets assets.Assets
}

func NewRenderSystem(
	world ecs.World,
	assets assets.Assets,
) RenderSystem {
	return RenderSystem{
		World:  world,
		Assets: assets,
	}
}

type renderableEntity struct {
	Id   ecs.EntityId
	Mesh mesh.MeshCachedAsset
}

type renderable struct {
	Material material.MaterialCachedAsset
	Entities []renderableEntity
}

func (s *RenderSystem) Update(args frames.FrameEvent) error {
	renderables := map[assets.AssetID]renderable{}

	renderableEntities := s.World.GetEntitiesWithComponents(
		ecs.GetComponentType(mesh.Mesh{}),
		ecs.GetComponentType(material.Material{}),
	)

	for _, entity := range renderableEntities {
		var materialComponent material.Material
		if err := s.World.GetComponents(entity, &materialComponent); err != nil {
			continue
		}
		for _, materialID := range materialComponent.IDs {
			materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](s.Assets, materialID)
			if err != nil {
				return err
			}

			var meshComponent mesh.Mesh
			if err := s.World.GetComponents(entity, &meshComponent); err != nil {
				continue
			}

			meshAsset, err := assets.GetAsset[mesh.MeshCachedAsset](s.Assets, meshComponent.ID)
			if err != nil {
				return err
			}

			if existing, ok := renderables[materialID]; ok {
				existing.Entities = append(existing.Entities, renderableEntity{entity, meshAsset})
				renderables[materialID] = existing
			} else {
				renderables[materialID] = renderable{
					materialAsset,
					[]renderableEntity{{entity, meshAsset}},
				}
			}
		}
	}

	for _, material := range renderables {
		if err := material.Material.OnFrame(s.World); err != nil {
			return err
		}

		for _, entity := range material.Entities {
			if err := material.Material.UseForEntity(s.World, entity.Id); err != nil {
				return err
			}
			entity.Mesh.VAO().Draw()
		}
	}

	return nil
}
