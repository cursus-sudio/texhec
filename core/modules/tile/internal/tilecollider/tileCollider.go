package tilecollider

import (
	"core/modules/tile"
	"engine"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"
)

func TileColliderSystem(
	w engine.World,
	tileService tile.Service,
	tileSize int32, // transform
	gridDepth float32,
	tileGroups groups.GroupsComponent, // groups
	colliderComponent collider.Component, // collider
) tile.System {
	return ecs.NewSystemRegister(func() error {
		tilePosDirtySet := ecs.NewDirtySet()
		tileService.Pos().AddDirtySet(tilePosDirtySet)
		w.UUID.Component().BeforeGet(func() {
			ei := tilePosDirtySet.Get()
			if len(ei) == 0 {
				return
			}
			for _, entity := range ei {
				if _, ok := w.UUID.Component().Get(entity); ok {
					continue
				}
				comp := uuid.New(w.UUID.NewUUID())
				w.UUID.Component().Set(entity, comp)
			}
		})

		//

		tileDirtySet := ecs.NewDirtySet()
		tileService.Pos().AddDirtySet(tileDirtySet)
		applyTileCollider := func() {
			ei := tileDirtySet.Get()
			if len(ei) == 0 {
				return
			}

			for _, entity := range ei {

				pos, ok := tileService.Pos().Get(entity)
				if !ok {
					continue
				}
				transformPos := transform.NewPos(
					float32(tileSize)*float32(pos.X)+float32(tileSize)/2,
					float32(tileSize)*float32(pos.Y)+float32(tileSize)/2,
					gridDepth+float32(pos.Layer),
				)
				w.Transform.Pos().Set(entity, transformPos)
				w.Transform.Size().Set(entity, transform.NewSize(float32(tileSize), float32(tileSize), 1))
				w.Inputs.LeftClick().Set(entity, inputs.NewLeftClick(tile.NewTileClickEvent(pos)))
				w.Inputs.Stack().Set(entity, inputs.StackComponent{})
				w.Collider.Component().Set(entity, colliderComponent)
				w.Groups.Component().Set(entity, tileGroups)
			}
		}

		w.Transform.Size().BeforeGet(applyTileCollider)
		w.Transform.Pos().BeforeGet(applyTileCollider)
		w.Collider.Component().BeforeGet(applyTileCollider)
		w.Inputs.LeftClick().BeforeGet(applyTileCollider)
		w.Inputs.Stack().BeforeGet(applyTileCollider)
		w.Groups.Component().BeforeGet(applyTileCollider)

		return nil
	})
}
