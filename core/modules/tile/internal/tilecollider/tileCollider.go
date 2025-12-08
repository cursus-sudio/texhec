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

		tilePosArray.OnAdd(func(ei []ecs.EntityID) {
			uuidTransaction := uuidArray.Transaction()
			for _, entity := range ei {
				if _, err := uuidArray.GetComponent(entity); err == nil {
					continue
				}
				comp := uuid.New(uuidFactory.NewUUID())
				uuidTransaction.SaveComponent(entity, comp)
			}
			logger.Warn(ecs.FlushMany(uuidTransaction))
		})

		onUpsert := func(ei []ecs.EntityID) {
			// groups
			groupsTransaction := groupsArray.Transaction()
			for _, entity := range ei {
				groupsTransaction.SaveComponent(entity, tileGroups)
			}

			// pos
			posTransaction := posArray.Transaction()
			leftClickTransaction := leftClickArray.Transaction()
			for _, entity := range ei {
				pos, err := tilePosArray.GetComponent(entity)
				if err != nil {
					continue
				}
				transformPos := transform.NewPos(
					float32(tileSize)*float32(pos.X)+float32(tileSize)/2,
					float32(tileSize)*float32(pos.Y)+float32(tileSize)/2,
					gridDepth+float32(pos.Layer),
				)
				posTransaction.SaveComponent(entity, transformPos)
				comp := inputs.NewMouseLeftClick(tile.NewTileClickEvent(pos))
				leftClickTransaction.SaveComponent(entity, comp)
			}

			// transform
			sizeTransaction := sizeArray.Transaction()
			for _, entity := range ei {
				sizeTransaction.SaveComponent(entity, transform.NewSize(float32(tileSize), float32(tileSize), 1))
			}

			// collider
			colliderTransaction := collidersArray.Transaction()
			for _, entity := range ei {
				colliderTransaction.SaveComponent(entity, colliderComponent)
			}

			// mouse
			logger.Warn(ecs.FlushMany(
				groupsTransaction,
				posTransaction,
				sizeTransaction,
				leftClickTransaction,
				colliderTransaction,
			))
		}

		tileColliderArray := ecs.GetComponentsArray[ColliderComponent](w)
		tileColliderArray.OnAdd(onUpsert)
		tileColliderArray.OnChange(onUpsert)
		return nil
	})
}
