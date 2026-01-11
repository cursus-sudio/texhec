package smoothpkg

import (
	"engine/modules/smooth"
	"engine/modules/smooth/internal"
	"engine/modules/transition"
	"engine/services/ecs"
	"reflect"

	"github.com/ogiusek/ioc/v2"
)

type config struct {
	components   map[reflect.Type]struct{}
	services     []func(b ioc.Builder)
	firstSystems []func(c ioc.Dic) smooth.StartSystem
	lastSystems  []func(c ioc.Dic) smooth.StopSystem
}

type Config struct {
	*config
}

func NewConfig() Config {
	return Config{
		config: &config{
			components: make(map[reflect.Type]struct{}),
		},
	}
}

type startSystem[Component any] smooth.StartSystem
type stopSystem[Component any] smooth.StopSystem

func SmoothComponent[Component transition.Lerp[Component]](config Config) {
	componentType := reflect.TypeFor[Component]()
	if _, ok := config.components[componentType]; ok {
		return
	}

	config.components[componentType] = struct{}{}
	config.services = append(config.services, func(b ioc.Builder) {
		ioc.RegisterSingleton(b, func(c ioc.Dic) *internal.Service[Component] {
			return internal.NewService[Component](c)
		})
		ioc.RegisterSingleton(b, func(c ioc.Dic) startSystem[Component] {
			return internal.NewFirstSystem[Component](c)
		})
		ioc.RegisterSingleton(b, func(c ioc.Dic) stopSystem[Component] {
			return internal.NewLastSystem[Component](c)
		})
	})
	config.firstSystems = append(config.firstSystems, func(c ioc.Dic) smooth.StartSystem {
		return ioc.Get[startSystem[Component]](c)
	})
	config.lastSystems = append(config.lastSystems, func(c ioc.Dic) smooth.StopSystem {
		return ioc.Get[stopSystem[Component]](c)
	})
}

type pkg struct {
	config Config
}

func Package(config Config) ioc.Pkg {
	return pkg{config}
}

func (pkg pkg) Register(b ioc.Builder) {
	for _, register := range pkg.config.services {
		register(b)
	}
	ioc.RegisterSingleton(b, func(c ioc.Dic) smooth.StartSystem {
		return ecs.NewSystemRegister(func() error {
			for _, system := range pkg.config.firstSystems {
				if err := system(c).Register(); err != nil {
					return err
				}
			}
			return nil
		})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) smooth.StopSystem {
		return ecs.NewSystemRegister(func() error {
			for _, system := range pkg.config.lastSystems {
				if err := system(c).Register(); err != nil {
					return err
				}
			}
			return nil
		})
	})
}
