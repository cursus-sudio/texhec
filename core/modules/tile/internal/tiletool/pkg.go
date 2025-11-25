package tiletool

import (
	"core/modules/tile"
	"engine/modules/relation"
	relationpkg "engine/modules/relation/pkg"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	tileSize                     int32
	gridDepth                    float32
	mainLayer                    tile.Layer
	layers                       []tile.Layer
	minX, maxX, minY, maxY, minZ int32
	relationPkgs                 []ioc.Pkg
}

func Package(
	tileSize int32,
	gridDepth float32,
	mainLayer tile.Layer,
	layers []tile.Layer,
	minX, maxX, minY, maxY, minZ int32,
) ioc.Pkg {
	return pkg{
		tileSize:  tileSize,
		gridDepth: gridDepth,
		mainLayer: mainLayer,
		layers:    layers,
		minX:      minX,
		maxX:      maxX,
		minY:      minY,
		maxY:      maxY,
		minZ:      minZ,
		relationPkgs: []ioc.Pkg{
			relationpkg.SpatialRelationPackage(
				func(w ecs.World) ecs.LiveQuery {
					return w.Query().
						Require(ecs.GetComponentType(tile.PosComponent{})).
						Build()
				},
				func(w ecs.World) func(entity ecs.EntityID) (tile.PosComponent, bool) {
					tilePosArray := ecs.GetComponentsArray[tile.PosComponent](w)
					return func(entity ecs.EntityID) (tile.PosComponent, bool) {
						comp, err := tilePosArray.GetComponent(entity)
						return comp, err == nil
					}
				},
				func(index tile.PosComponent) uint32 {
					xMul := maxX - minX
					yMul := xMul * (maxY - minY)
					result := (index.X+minX)*xMul + (index.Y+minY)*yMul + (int32(index.Layer) + minZ)
					return uint32(result)
				},
			),
			relationpkg.SpatialRelationPackage(
				func(w ecs.World) ecs.LiveQuery {
					return w.Query().
						Require(ecs.GetComponentType(tile.PosComponent{})).
						Build()
				},
				func(w ecs.World) func(entity ecs.EntityID) (tile.ColliderPos, bool) {
					tilePosArray := ecs.GetComponentsArray[tile.PosComponent](w)
					return func(entity ecs.EntityID) (tile.ColliderPos, bool) {
						tileComp, err := tilePosArray.GetComponent(entity)
						if err != nil && tileComp.Layer != mainLayer {
							return tile.ColliderPos{}, false
						}
						return tileComp.GetColliderPos(), true
					}
				},
				func(pos tile.ColliderPos) uint32 {
					xMul := maxX - minX
					result := (pos.X+minX)*xMul + (pos.Y + minY)
					return uint32(result)
				},
			),
		},
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.relationPkgs {
		pkg.Register(b)
	}
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[tile.Tool] {
		return ecs.NewToolFactory(func(w ecs.World) tile.Tool {
			return &tool{
				ioc.Get[ecs.ToolFactory[relation.EntityToKeyTool[tile.PosComponent]]](c).Build(w),
				ioc.Get[ecs.ToolFactory[relation.EntityToKeyTool[tile.ColliderPos]]](c).Build(w),
			}
		})
	})
}
