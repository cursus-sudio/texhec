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
	firstSystems []smooth.StartSystem
	lastSystems  []smooth.StopSystem
}

type Config struct {
	*config
}

func NewConfig() Config {
	return Config{
		config: &config{
			components:   make(map[reflect.Type]struct{}),
			firstSystems: make([]smooth.StartSystem, 0),
			lastSystems:  make([]smooth.StopSystem, 0),
		},
	}
}

func SmoothComponent[Component transition.Lerp[Component]](config Config) {
	componentType := reflect.TypeFor[Component]()
	if _, ok := config.components[componentType]; ok {
		return
	}

	config.components[componentType] = struct{}{}
	config.firstSystems = append(config.firstSystems, internal.NewFirstSystem[Component]())
	config.lastSystems = append(config.lastSystems, internal.NewLastSystem[Component]())
}

type pkg struct {
	config Config
}

func Package(config Config) ioc.Pkg {
	return pkg{config}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) smooth.StartSystem {
		return ecs.NewSystemRegister(func(w smooth.World) error {
			for _, system := range pkg.config.firstSystems {
				if err := system.Register(w); err != nil {
					return err
				}
			}
			return nil
		})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) smooth.StopSystem {
		return ecs.NewSystemRegister(func(w smooth.World) error {
			for _, system := range pkg.config.lastSystems {
				if err := system.Register(w); err != nil {
					return err
				}
			}
			return nil
		})
	})
}
