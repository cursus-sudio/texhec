package test

import (
	"engine/modules/relation"
	"engine/modules/relation/pkg"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type Component struct {
	Index uint32
}

type Setup struct {
	W     ecs.World
	Array ecs.ComponentsArray[Component]
	Tool  func() relation.EntityToKeyTool[uint32]
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	pkgs := []ioc.Pkg{
		relationpkg.SpatialRelationPackage(
			func(w ecs.World) ecs.DirtySet {
				dirtySet := ecs.NewDirtySet()
				ecs.GetComponentsArray[Component](w).AddDirtySet(dirtySet)
				return dirtySet
			},
			func(w ecs.World) func(entity ecs.EntityID) (uint32, bool) {
				componentArray := ecs.GetComponentsArray[Component](w)
				return func(entity ecs.EntityID) (uint32, bool) {
					comp, ok := componentArray.Get(entity)
					return comp.Index, ok
				}
			},
			func(index uint32) uint32 { return index },
		),
	}
	for _, pkg := range pkgs {
		pkg.Register(b)
	}

	c := b.Build()
	toolFactory := ioc.Get[ecs.ToolFactory[relation.EntityToKeyTool[uint32]]](c)

	w := ecs.NewWorld()
	return Setup{
		W:     w,
		Array: ecs.GetComponentsArray[Component](w),
		Tool:  func() relation.EntityToKeyTool[uint32] { return toolFactory.Build(w) },
	}
}
