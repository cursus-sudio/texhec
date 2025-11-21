package parent_test

import (
	"frontend/modules/relation"
	"frontend/modules/relation/pkg"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type Component struct {
	Parent ecs.EntityID
}

type Setup struct {
	W     ecs.World
	Array ecs.ComponentsArray[Component]
	Tool  func() relation.ParentTool[Component]
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	pkgs := []ioc.Pkg{
		relationpkg.ParentPackage(func(c Component) ecs.EntityID { return c.Parent }),
	}
	for _, pkg := range pkgs {
		pkg.Register(b)
	}

	c := b.Build()
	toolFactory := ioc.Get[ecs.ToolFactory[relation.ParentTool[Component]]](c)

	w := ecs.NewWorld()
	return Setup{
		W:     w,
		Array: ecs.GetComponentsArray[Component](w),
		Tool:  func() relation.ParentTool[Component] { return toolFactory.Build(w) },
	}
}
