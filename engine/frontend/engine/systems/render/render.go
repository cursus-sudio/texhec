package render

import (
	"frontend/engine/components/material"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/frames"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type RenderSystem struct {
	world  ecs.World
	assets assets.Assets

	setFence bool
	fence    uintptr

	materials map[assets.AssetID]material.MaterialCachedAsset
}

func NewRenderSystem(
	world ecs.World,
	assetsService assets.Assets,
) RenderSystem {
	liveMaterials := map[assets.AssetID]material.MaterialCachedAsset{}

	liveQuery := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(material.Material{}),
	)

	onChange := func(_ []ecs.EntityID) {
		materialAssets := map[assets.AssetID]struct{}{}
		for _, entity := range liveQuery.Entities() {
			materialComponent, err := ecs.GetComponent[material.Material](world, entity)
			if err != nil {
				continue
			}
			for _, id := range materialComponent.IDs {
				materialAssets[id] = struct{}{}
			}
		}

		for assetID := range liveMaterials {
			if _, ok := materialAssets[assetID]; !ok {
				delete(liveMaterials, assetID)
			}
		}

		for assetID := range materialAssets {
			if _, ok := liveMaterials[assetID]; ok {
				continue
			}

			materialAsset, err := assets.GetAsset[material.MaterialCachedAsset](assetsService, assetID)
			if err != nil {
				delete(liveMaterials, assetID)
				continue
			}

			liveMaterials[assetID] = materialAsset
		}

	}

	liveQuery.OnAdd(onChange)
	liveQuery.OnRemove(onChange)

	return RenderSystem{
		world:  world,
		assets: assetsService,

		materials: liveMaterials}
}

type renderable struct {
	Material material.MaterialCachedAsset
}

func (s *RenderSystem) Listen(args frames.FrameEvent) error {
	if s.setFence {
		s.setFence = false
		gl.ClientWaitSync(s.fence, gl.SYNC_FLUSH_COMMANDS_BIT, gl.TIMEOUT_IGNORED)
		gl.DeleteSync(s.fence)
	}
	for _, material := range s.materials {
		if err := material.Render(s.world); err != nil {
			return err
		}
	}

	s.fence = gl.FenceSync(gl.SYNC_GPU_COMMANDS_COMPLETE, 0)
	s.setFence = true

	return nil
}
