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
	Tool  func() indexing.Indices[uint32]
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	pkgs := []ioc.Pkg{
		indexingpkg.SpatialIndexPackage(
			func(w ecs.World) ecs.LiveQuery {
				return w.Query().
					Require(ecs.GetComponentType(Component{})).
					Build()
			},
			func(w ecs.World) func(entity ecs.EntityID) (uint32, bool) {
				componentArray := ecs.GetComponentsArray[Component](w)
				return func(entity ecs.EntityID) (uint32, bool) {
					comp, err := componentArray.GetComponent(entity)
					return comp.Index, err == nil
				}
			},
			func(index uint32) uint32 { return index },
		),
	}
	for _, pkg := range pkgs {
		pkg.Register(b)
	}

	c := b.Build()
	toolFactory := ioc.Get[ecs.ToolFactory[indexing.Indices[uint32]]](c)

	w := ecs.NewWorld()
	return Setup{
		W:     w,
		Array: ecs.GetComponentsArray[Component](w),
		Tool:  func() indexing.Indices[uint32] { return toolFactory.Build(w) },
	}
}
