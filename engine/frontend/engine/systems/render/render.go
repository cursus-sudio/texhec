package render

import (
	"frontend/engine/components/material"
	"frontend/engine/components/mesh"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/frames"
	"shared/services/logger"
)

type RenderSystem struct {
	World  ecs.World
	Assets assets.Assets
	Logger logger.Logger
}

func NewRenderSystem(world ecs.World, assets assets.Assets, logger logger.Logger) RenderSystem {
	return RenderSystem{
		World:  world,
		Assets: assets,
		Logger: logger,
	}
}

func (s *RenderSystem) Update(args frames.FrameEvent) {
	materials := map[assets.AssetID]material.MaterialCachedAsset{}
	for _, entity := range s.World.GetEntitiesWithComponents(
		ecs.GetComponentType(material.Material{}),
	) {
		var materialComponent material.Material
		if err := s.World.GetComponent(entity, &materialComponent); err != nil {
			continue
		}
		materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](s.Assets, materialComponent.ID)
		if err != nil {
			s.Logger.Error(err)
			continue
		}
		materials[materialComponent.ID] = materialAsset
	}

	for _, material := range materials {
		if err := material.OnFrame(s.World); err != nil {
			s.Logger.Error(err)
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
			s.Logger.Error(err)
			continue
		}

		var materialComponent material.Material
		if err := s.World.GetComponent(entity, &materialComponent); err != nil {
			continue
		}
		materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](s.Assets, materialComponent.ID)
		if err != nil {
			s.Logger.Error(err)
			continue
		}

		if err := materialAsset.UseForEntity(s.World, entity); err != nil {
			s.Logger.Error(err)
			continue
		}

		meshAsset.VAO().Draw()
	}
}
