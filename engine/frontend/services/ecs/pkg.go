package ecs

import (
	"frontend/services/ecs/ecsargs"
	"frontend/services/frames"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) WorldFactory { return func() World { return newWorld() } })

	ioc.WrapService(b, frames.Update, func(c ioc.Dic, f frames.Builder) frames.Builder {
		return f.OnFrame(func(of frames.OnFrame) {
			world := ioc.Get[World](c)
			deltaTime := ecsargs.NewDeltaTime(of.Delta)
			args := NewArgs(deltaTime)
			world.Update(args)
		})
	})
}
