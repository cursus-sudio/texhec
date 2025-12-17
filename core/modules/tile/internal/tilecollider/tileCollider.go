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
	colliderComponent collider.ColliderComponent, // collider
	uuidFactory uuid.Factory, // tools
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		uuidArray := ecs.GetComponentsArray[uuid.Component](w)
		tilePosArray := ecs.GetComponentsArray[tile.PosComponent](w)

		leftClickArray := ecs.GetComponentsArray[inputs.MouseLeftClickComponent](w)
		collidersArray := ecs.GetComponentsArray[collider.ColliderComponent](w)

		posArray := ecs.GetComponentsArray[transform.PosComponent](w)
		sizeArray := ecs.GetComponentsArray[transform.SizeComponent](w)

		groupsArray := ecs.GetComponentsArray[groups.GroupsComponent](w)

		tilePosDirtySet := ecs.NewDirtySet()
		tilePosArray.AddDirtySet(tilePosDirtySet)
		uuidArray.BeforeGet(func() {
			ei := tilePosDirtySet.Get()
			if len(ei) == 0 {
				return
			}
			for _, entity := range ei {
				if _, ok := uuidArray.Get(entity); ok {
					continue
				}
				comp := uuid.New(uuidFactory.NewUUID())
				uuidArray.Set(entity, comp)
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
				groupsArray.Set(entity, tileGroups)
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
				posArray.Set(entity, transformPos)
				comp := inputs.NewMouseLeftClick(tile.NewTileClickEvent(pos))
				leftClickArray.Set(entity, comp)
			}

			// transform
			for _, entity := range ei {
				sizeArray.Set(entity, transform.NewSize(float32(tileSize), float32(tileSize), 1))
			}

			// collider
			for _, entity := range ei {
				collidersArray.Set(entity, colliderComponent)
			}
		}

		collidersArray.BeforeGet(applyTileCollider)
		sizeArray.BeforeGet(applyTileCollider)
		posArray.BeforeGet(applyTileCollider)
		leftClickArray.BeforeGet(applyTileCollider)
		groupsArray.BeforeGet(applyTileCollider)

		return nil
	})
}
