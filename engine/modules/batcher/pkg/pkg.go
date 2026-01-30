package batcherpkg

import (
	"engine/modules/batcher"
	"engine/modules/batcher/internal"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	workers            int
	frameLoadingBudget time.Duration
}

func Package(
	workers int,
	frameLoadingBudget time.Duration,
) ioc.Pkg {
	return pkg{
		workers,
		frameLoadingBudget,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) *internal.Service {
		return internal.NewService(
			c,
			pkg.workers,
			pkg.frameLoadingBudget,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) batcher.Service {
		return ioc.Get[*internal.Service](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) batcher.System {
		return ioc.Get[*internal.Service](c).System()
	})
}
