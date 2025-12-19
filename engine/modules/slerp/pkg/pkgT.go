package slerppkg

import (
	"engine/modules/slerp"
	"engine/modules/slerp/internal/sys"

	"github.com/ogiusek/ioc/v2"
)

type pkgT[Component any] struct {
	slerpFn sys.SlerpFn[Component]
}

func PackageT[Component any](
	slerpFn func(c1, c2 Component, progress slerp.Progress) Component,
) ioc.Pkg {
	return pkgT[Component]{slerpFn}
}

func (pkg pkgT[Component]) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b sys.Builder) sys.Builder {
		sys := sys.NewSysT(pkg.slerpFn)
		b.Register(sys)
		return b
	})
}
