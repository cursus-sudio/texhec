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

func TileColliderSystem(logger logger.Logger,
	tileSize int32, // transform
	gridDepth float32,
	tileGroups groups.GroupsComponent, // groups
	colliderComponent collider.Component, // collider
	uuidFactory uuid.Factory, // tools
) ecs.SystemRegister[tile.World] {
	return ecs.NewSystemRegister(func(w tile.World) error {
		tilePosArray := ecs.GetComponentsArray[tile.PosComponent](w)

		tilePosDirtySet := ecs.NewDirtySet()
		tilePosArray.AddDirtySet(tilePosDirtySet)
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

		tileColliderDirtySet := ecs.NewDirtySet()
		tileColliderArray := ecs.GetComponentsArray[ColliderComponent](w)
		tileColliderArray.AddDirtySet(tileColliderDirtySet)
		applyTileCollider := func() {
			ei := tileColliderDirtySet.Get()
			if len(ei) == 0 {
				return
			}
			// groups
			for _, entity := range ei {
				w.Groups().Component().Set(entity, tileGroups)
			}

			// pos
			for _, entity := range ei {
				pos, ok := tilePosArray.Get(entity)
				if !ok {
					continue
				}
				transformPos := transform.NewPos(
					float32(tileSize)*float32(pos.X)+float32(tileSize)/2,
					float32(tileSize)*float32(pos.Y)+float32(tileSize)/2,
					gridDepth+float32(pos.Layer),
				)
				w.Transform().Pos().Set(entity, transformPos)
				comp := inputs.NewMouseLeftClick(tile.NewTileClickEvent(pos))
				w.Inputs().MouseLeft().Set(entity, comp)
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
		w.Inputs().MouseLeft().BeforeGet(applyTileCollider)
		w.Groups().Component().BeforeGet(applyTileCollider)

		return nil
	})
}
