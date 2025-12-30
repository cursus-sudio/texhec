package tilecollider

import (
	"core/modules/tile"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"
)

func TileColliderSystem(
	toolFactory tile.ToolFactory,
	logger logger.Logger,
	tileSize int32, // transform
	gridDepth float32,
	tileGroups groups.GroupsComponent, // groups
	colliderComponent collider.Component, // collider
	uuidFactory uuid.Factory, // tools
) tile.System {
	return ecs.NewSystemRegister(func(w tile.World) error {
		tileTool := toolFactory.Build(w)
		tilePosDirtySet := ecs.NewDirtySet()
		tileTool.Tile().Pos().AddDirtySet(tilePosDirtySet)
		w.UUID().Component().BeforeGet(func() {
			ei := tilePosDirtySet.Get()
			if len(ei) == 0 {
				return
			}
			for _, entity := range ei {
				if _, ok := w.UUID().Component().Get(entity); ok {
					continue
				}
				comp := uuid.New(uuidFactory.NewUUID())
				w.UUID().Component().Set(entity, comp)
			}
		})

		//

		tileDirtySet := ecs.NewDirtySet()
		tileTool.Tile().Pos().AddDirtySet(tileDirtySet)
		applyTileCollider := func() {
			ei := tileDirtySet.Get()
			if len(ei) == 0 {
				return
			}
			// groups
			for _, entity := range ei {
				w.Groups().Component().Set(entity, tileGroups)
			}

			// pos
			for _, entity := range ei {
				pos, ok := tileTool.Tile().Pos().Get(entity)
				if !ok {
					continue
				}
				transformPos := transform.NewPos(
					float32(tileSize)*float32(pos.X)+float32(tileSize)/2,
					float32(tileSize)*float32(pos.Y)+float32(tileSize)/2,
					gridDepth+float32(pos.Layer),
				)
				w.Transform().Pos().Set(entity, transformPos)
				comp := inputs.NewLeftClick(tile.NewTileClickEvent(pos))
				w.Inputs().LeftClick().Set(entity, comp)
				w.Inputs().Stack().Set(entity, inputs.StackComponent{})
			}

			// transform
			for _, entity := range ei {
				w.Transform().Size().Set(entity, transform.NewSize(float32(tileSize), float32(tileSize), 1))
			}

			// collider
			for _, entity := range ei {
				w.Collider().Component().Set(entity, colliderComponent)
			}
		}

		w.Collider().Component().BeforeGet(applyTileCollider)
		w.Transform().Size().BeforeGet(applyTileCollider)
		w.Transform().Pos().BeforeGet(applyTileCollider)
		w.Inputs().LeftClick().BeforeGet(applyTileCollider)
		w.Inputs().Stack().BeforeGet(applyTileCollider)
		w.Groups().Component().BeforeGet(applyTileCollider)

		return nil
	})
}
