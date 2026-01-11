package transformpkg

import (
	"engine/modules/transform"
	"engine/modules/transform/internal/transformservice"
	transitionpkg "engine/modules/transition/pkg"
	"engine/services/codec"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	defaultPos         transform.PosComponent
	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent
}

func Package() ioc.Pkg {
	return pkg{
		transform.NewPos(0, 0, 0),
		transform.NewRotation(mgl32.QuatIdent()),
		transform.NewSize(1, 1, 1),
		transform.NewPivotPoint(.5, .5, .5),
		transform.NewParentPivotPoint(.5, .5, .5),
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(transform.PosComponent{}).
			Register(transform.RotationComponent{}).
			Register(transform.SizeComponent{}).
			Register(transform.PivotPointComponent{}).
			Register(transform.ParentPivotPointComponent{}).
			Register(transform.ParentComponent{})
	})

	for _, pkg := range []ioc.Pkg{
		transitionpkg.PackageT[transform.PosComponent](),
		transitionpkg.PackageT[transform.RotationComponent](),
		transitionpkg.PackageT[transform.SizeComponent](),
		transitionpkg.PackageT[transform.PivotPointComponent](),
		transitionpkg.PackageT[transform.ParentPivotPointComponent](),
	} {
		pkg.Register(b)
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) transform.Service {
		return transformservice.NewService(c,
			pkg.defaultRot,
			pkg.defaultSize,
			pkg.defaultPivot,
			pkg.defaultParentPivot,
		)
	})

}
