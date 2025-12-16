package tilecollider

import (
	gameassets "core/assets"
	"core/modules/tile"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	tileSize                     int32
	gridDepth                    float32
	tileGroups                   groups.GroupsComponent
	mainLayer                    tile.Layer
	layers                       []tile.Layer
	minX, maxX, minY, maxY, minZ int32
}

func Package(
	tileSize int32,
	gridDepth float32,
	tileGroups groups.GroupsComponent,
	mainLayer tile.Layer,
	layers []tile.Layer,
	minX, maxX, minY, maxY, minZ int32,
) ioc.Pkg {
	return pkg{
		tileSize,
		gridDepth,
		tileGroups,
		mainLayer,
		layers,
		minX, maxX, minY, maxY, minZ,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.System) tile.System {
		tileToolFactory := ioc.Get[ecs.ToolFactory[tile.Tile]](c)
		logger := ioc.Get[logger.Logger](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			if err := s.Register(w); err != nil {
				return err
			}
			posIndex := tileToolFactory.Build(w).Tile().TilePos()
			errs := ecs.RegisterSystems(w,
				TileColliderSystem(
					logger,
					pkg.tileSize,
					pkg.gridDepth,
					pkg.tileGroups,
					collider.NewCollider(ioc.Get[gameassets.GameAssets](c).SquareCollider),
					ioc.Get[uuid.Factory](c),
				),
				ecs.NewSystemRegister(func(w ecs.World) error {
					entitiesPositions := datastructures.NewSparseArray[ecs.EntityID, tile.PosComponent]()
					dirtyEntities := ecs.NewDirtySet()

					posArray := ecs.GetComponentsArray[tile.PosComponent](w)
					tileColliderArray := ecs.GetComponentsArray[ColliderComponent](w)
					colliderArray := ecs.GetComponentsArray[collider.ColliderComponent](w)

					posArray.AddDirtySet(dirtyEntities)

					colliderArray.BeforeGet(func() {
						entities := dirtyEntities.Get()
						if len(entities) == 0 {
							return
						}
						finalColliders := datastructures.NewSparseArray[ecs.EntityID, ColliderComponent]()
						for _, entity := range entities {
							if comp, ok := entitiesPositions.Get(entity); ok {
								key := comp
								key.Layer = pkg.mainLayer
								if colliderEntity, ok := posIndex.Get(key); ok {
									collider, ok := finalColliders.Get(colliderEntity)
									if !ok {
										collider, ok = tileColliderArray.GetComponent(colliderEntity)
									}
									if !ok {
										collider = NewCollider()
										collider.Add(pkg.mainLayer)
									}
									collider.LayersBitmask = collider.LayersBitmask &^ uint8(comp.Layer)
									finalColliders.Set(colliderEntity, collider)
								}
							}

							pos, ok := posArray.GetComponent(entity)
							if !ok {
								continue
							}
							key := pos
							key.Layer = pkg.mainLayer
							if colliderEntity, ok := posIndex.Get(key); ok {
								collider, ok := finalColliders.Get(colliderEntity)
								if !ok {
									collider, ok = tileColliderArray.GetComponent(colliderEntity)
								}
								if !ok {
									collider = NewCollider()
									collider.Add(pkg.mainLayer)
								}
								collider.LayersBitmask = collider.LayersBitmask | uint8(pos.Layer)
								finalColliders.Set(colliderEntity, collider)
							}
						}

						for _, entity := range finalColliders.GetIndices() {
							value, ok := finalColliders.Get(entity)
							if !ok {
								continue
							}
							tileColliderArray.SaveComponent(entity, value)
						}
					})
					return nil
				}),
			)
			if len(errs) != 0 {
				return errs[0]
			}
			return nil
		})
	})
}
