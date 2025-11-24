package tilecollider

import (
	"core/modules/tile"
	"frontend/modules/collider"
	"frontend/modules/groups"
	"frontend/modules/relation"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	tileSize                     int32
	gridDepth                    float32
	tileGroups                   groups.GroupsComponent
	colliderComponent            collider.ColliderComponent
	mainLayer                    tile.Layer
	layers                       []tile.Layer
	minX, maxX, minY, maxY, minZ int32
}

func Package(
	tileSize int32,
	gridDepth float32,
	tileGroups groups.GroupsComponent,
	colliderComponent collider.ColliderComponent,
	mainLayer tile.Layer,
	layers []tile.Layer,
	minX, maxX, minY, maxY, minZ int32,
) ioc.Pkg {
	return pkg{
		tileSize,
		gridDepth,
		tileGroups,
		colliderComponent,
		mainLayer,
		layers,
		minX, maxX, minY, maxY, minZ,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.System) tile.System {
		tileToolFactory := ioc.Get[ecs.ToolFactory[tile.Tool]](c)
		logger := ioc.Get[logger.Logger](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			if err := s.Register(w); err != nil {
				return err
			}
			posIndex := tileToolFactory.Build(w).TilePos()
			errs := ecs.RegisterSystems(w,
				TileColliderSystem(
					logger,
					pkg.tileSize,
					pkg.gridDepth,
					pkg.tileGroups,
					pkg.colliderComponent,
				),
				ecs.NewSystemRegister(func(w ecs.World) error {
					posArray := ecs.GetComponentsArray[tile.PosComponent](w)
					colliderArray := ecs.GetComponentsArray[ColliderComponent](w)
					posArray.OnRemoveComponents(func(ei []ecs.EntityID, components []tile.PosComponent) {
						colliderTransaction := colliderArray.Transaction()
						set := datastructures.NewSparseSet[ecs.EntityID]()
						for _, component := range components {
							component.Layer = pkg.mainLayer
							entity, ok := posIndex.Get(component)
							if ok {
								set.Add(entity)
							}
						}
						for _, entity := range set.GetIndices() {
							pos, err := posArray.GetComponent(entity)
							if err != nil {
								continue
							}
							collider := NewCollider().Ptr().Add(pkg.mainLayer).Val()
							for _, layer := range pkg.layers {
								pos.Layer = layer
								_, ok := posIndex.Get(pos)
								if ok {
									collider.Add(layer)
								}
							}
							colliderTransaction.SaveComponent(entity, collider)
						}
						logger.Warn(ecs.FlushMany(colliderTransaction))
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

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, indices ecs.ToolFactory[relation.EntityToKeyTool[tile.ColliderPos]]) ecs.ToolFactory[relation.EntityToKeyTool[tile.ColliderPos]] {
		posIndexFactory := ioc.Get[ecs.ToolFactory[relation.EntityToKeyTool[tile.PosComponent]]](c)
		logger := ioc.Get[logger.Logger](c)

		return ecs.NewToolFactory(func(w ecs.World) relation.EntityToKeyTool[tile.ColliderPos] {
			posIndex := posIndexFactory.Build(w)
			posArray := ecs.GetComponentsArray[tile.PosComponent](w)
			colliderArray := ecs.GetComponentsArray[ColliderComponent](w)
			upsertEntities := func(ei []ecs.EntityID) {
				colliderTransaction := colliderArray.Transaction()
				for _, entity := range ei {
					pos, err := posArray.GetComponent(entity)
					if err != nil {
						continue
					}
					collider := NewCollider().Ptr().Add(pkg.mainLayer).Val()
					for _, layer := range pkg.layers {
						pos.Layer = layer
						_, ok := posIndex.Get(pos)
						if ok {
							collider.Add(layer)
						}
					}
					colliderTransaction.SaveComponent(entity, collider)
				}
				logger.Warn(ecs.FlushMany(colliderTransaction))
			}
			tool := indices.Build(w)
			tool.OnUpsert(upsertEntities)
			return tool
		})
	})
}
