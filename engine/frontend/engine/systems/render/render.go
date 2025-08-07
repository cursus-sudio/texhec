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

func (s *RenderSystem) Update(args frames.FrameEvent) error {
	materials := map[assets.AssetID]material.MaterialCachedAsset{}
	for _, entity := range s.World.GetEntitiesWithComponents(
		ecs.GetComponentType(material.Material{}),
	) {
		var materialComponent material.Material
		if err := s.World.GetComponent(entity, &materialComponent); err != nil {
			continue
		}
		for _, materialID := range materialComponent.IDs {
			materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](s.Assets, materialID)
			if err != nil {
				return err
			}
			materials[materialID] = materialAsset
		}
	}

	for _, material := range materials {
		if err := material.OnFrame(s.World); err != nil {
			return err
		}
	}

	renderableEntities := s.World.GetEntitiesWithComponents(
		ecs.GetComponentType(mesh.Mesh{}),
		ecs.GetComponentType(material.Material{}),
	)
	for _, entity := range renderableEntities {
		var meshComponent mesh.Mesh
		if err := s.World.GetComponent(entity, &meshComponent); err != nil {
			continue
		}
		meshAsset, err := assets.GetAsset[mesh.MeshCachedAsset](s.Assets, meshComponent.ID)
		if err != nil {
			return err
		}

		var materialComponent material.Material
		if err := s.World.GetComponent(entity, &materialComponent); err != nil {
			continue
		}
		for _, materialID := range materialComponent.IDs {
			materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](s.Assets, materialID)
			if err != nil {
				return err
			}

			if err := materialAsset.UseForEntity(s.World, entity); err != nil {
				return err
			}

			meshAsset.VAO().Draw()
		}
	}

	return nil
}
