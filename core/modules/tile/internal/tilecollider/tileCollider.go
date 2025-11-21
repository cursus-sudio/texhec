package tilecollider

import (
	"core/modules/tile"
	"frontend/modules/collider"
	"frontend/modules/groups"
	"frontend/modules/inputs"
	"frontend/modules/transform"
	"shared/services/ecs"
	"shared/services/logger"
)

func TileColliderSystem(
	logger logger.Logger,
	// transform
	tileSize int32,
	gridDepth float32,
	// groups
	tileGroups groups.GroupsComponent,
	// collider
	colliderComponent collider.ColliderComponent,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		tilePosArray := ecs.GetComponentsArray[tile.PosComponent](w.Components())
		tileColliderArray := ecs.GetComponentsArray[ColliderComponent](w.Components())

		leftClickArray := ecs.GetComponentsArray[inputs.MouseLeftClickComponent](w.Components())
		collidersArray := ecs.GetComponentsArray[collider.ColliderComponent](w.Components())

		posArray := ecs.GetComponentsArray[transform.PosComponent](w.Components())
		sizeArray := ecs.GetComponentsArray[transform.SizeComponent](w.Components())

		groupsArray := ecs.GetComponentsArray[groups.GroupsComponent](w.Components())

		onUpsert := func(ei []ecs.EntityID) {
			// groups
			groupsTransaction := groupsArray.Transaction()
			for _, entity := range ei {
				groupsTransaction.SaveComponent(entity, tileGroups)
			}

			// pos
			posTransaction := posArray.Transaction()
			for _, entity := range ei {
				pos, err := tilePosArray.GetComponent(entity)
				if err != nil {
					continue
				}
				posTransaction.SaveComponent(entity, transform.NewPos(
					float32(tileSize)*float32(pos.X)+float32(tileSize)/2,
					float32(tileSize)*float32(pos.Y)+float32(tileSize)/2,
					gridDepth+float32(pos.Layer),
				))
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
			leftClickTransaction := leftClickArray.Transaction()
			for _, entity := range ei {
				comp := inputs.NewMouseLeftClick(tile.NewTileClickEvent(entity))
				leftClickTransaction.SaveComponent(entity, comp)
			}
			logger.Warn(ecs.FlushMany(
				groupsTransaction,
				posTransaction,
				sizeTransaction,
				leftClickTransaction,
				colliderTransaction,
			))
		}

		tileColliderArray.OnAdd(onUpsert)
		tileColliderArray.OnChange(onUpsert)
		return nil
	})
}
