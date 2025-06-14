package scenes

import (
	"frontend/src/engine/ecs"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterTransient(c, func(c ioc.Dic) ecs.World { return ioc.Get[SceneManager](c).CurrentScene().World() })
	ioc.RegisterSingleton(c, func(c ioc.Dic) SceneManager { return newSceneManager() })
	ioc.RegisterSingleton(c, func(c ioc.Dic) SceneBuilder { return newSceneBuilder() })
}
