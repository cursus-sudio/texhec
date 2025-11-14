package test

import (
	"frontend/modules/indexing"
	indexingpkg "frontend/modules/indexing/pkg"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type Component struct {
	Index uint32
}

type Setup struct {
	W     ecs.World
	Array ecs.ComponentsArray[Component]
	Tool  func() indexing.SpatialIndexTool[uint32]
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	pkgs := []ioc.Pkg{
		indexingpkg.SpatialIndexingPackage(
			func(w ecs.World) ecs.LiveQuery {
				return w.Query().
					Require(ecs.GetComponentType(Component{})).
					Build()
			},
			func(w ecs.World) func(entity ecs.EntityID) uint32 {
				componentArray := ecs.GetComponentsArray[Component](w.Components())
				return func(entity ecs.EntityID) uint32 {
					comp, err := componentArray.GetComponent(entity)
					if err != nil {
						return 0
					}
					return comp.Index
				}
			},
			func(index uint32) uint32 { return index },
		),
	}
	for _, pkg := range pkgs {
		pkg.Register(b)
	}

	c := b.Build()
	toolFactory := ioc.Get[ecs.ToolFactory[indexing.SpatialIndexTool[uint32]]](c)

	w := ecs.NewWorld()
	return Setup{
		W:     w,
		Array: ecs.GetComponentsArray[Component](w.Components()),
		Tool:  func() indexing.SpatialIndexTool[uint32] { return toolFactory.Build(w) },
	}
}
